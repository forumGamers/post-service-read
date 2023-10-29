package documents

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"time"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	i "github.com/forumGamers/post-service-read/interfaces"
	v "github.com/forumGamers/post-service-read/validator"
	"github.com/forumGamers/post-service-read/web"
	"github.com/olivere/elastic/v7"
)

type PostService interface {
	Insert(ctx context.Context, data PostDocument) error
	FindById(ctx context.Context, id string) (json.RawMessage, error)
	DeleteOneById(ctx context.Context, id string) error
	GetPublicContent(ctx context.Context, query web.PostParams) ([]i.PostResponse, struct {
		Total    int
		Relation string
	}, error)
	BulkCreate(ctx context.Context, datas []PostDocument) error
}

type Media struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PostDocument struct {
	Id           string `json:"id"`
	UserId       string `json:"userId"`
	Text         string `json:"text" bson:"text"`
	Media        Media
	AllowComment bool `json:"allowComment" default:"true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tags         []string `json:"tags"`
	Privacy      string   `json:"privacy" default:"Public"`
}

func (p *PostDocument) Insert(ctx context.Context, data PostDocument) error {
	return database.NewIndex(database.POSTINDEX).InsertOne(ctx, data)
}

func (p *PostDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.POSTINDEX).DeleteOne(ctx, id)
}

func NewPost() PostService {
	return &PostDocument{}
}

func (p *PostDocument) FindById(ctx context.Context, id string) (json.RawMessage, error) {
	get, err := database.DB.Get().Index(database.POSTINDEX).Id(id).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, h.NotFound
		}

		return nil, err
	}

	return get.Source, nil
}

func (p *PostDocument) GetPublicContent(ctx context.Context, query web.PostParams) ([]i.PostResponse, struct {
	Total    int
	Relation string
}, error) {
	search := database.DB.Search().
		Index(database.POSTINDEX).
		Size(query.Limit).
		Sort("CreatedAt", false).
		Sort("id", false)

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewRangeQuery("CreatedAt").Gte("now-3d/d").Lte("now/d"))

	if len(query.UserIds) > 0 {
		var ids []any
		for _, val := range query.UserIds {
			ids = append(ids, val)
		}
		boolQuery.Must(elastic.NewTermsQuery("userId", ids...))
	}
	search.Query(boolQuery)

	if query.Page != nil {
		var sa []any
		timeStamp, err := strconv.ParseInt(query.Page[0], 10, 64)
		if err != nil && regexp.MustCompile(v.RegexID).MatchString(query.Page[1]) {
			sa = append(sa, timeStamp, query.Page[1])
			search.SearchAfter(sa...)
		}
	}

	result, err := search.Do(context.Background())
	if err != nil {
		return nil, struct {
			Total    int
			Relation string
		}{}, err
	}

	var postResponses []i.PostResponse
	if result.Hits.TotalHits.Value > 0 {
		for _, hit := range result.Hits.Hits {
			var post PostDocument
			json.Unmarshal(hit.Source, &post)
			postResponses = append(postResponses, i.PostResponse{
				Id:           post.Id,
				UserId:       post.UserId,
				Text:         post.Text,
				Media:        i.Media(post.Media),
				AllowComment: post.AllowComment,
				CreatedAt:    post.CreatedAt,
				UpdatedAt:    post.UpdatedAt,
				Tags:         post.Tags,
				Privacy:      post.Privacy,
				SearchAfter:  hit.Sort,
			})
		}
	}

	return postResponses, struct {
		Total    int
		Relation string
	}{int(result.Hits.TotalHits.Value), result.Hits.TotalHits.Relation}, nil
}

func (p *PostDocument) BulkCreate(ctx context.Context, datas []PostDocument) error {
	bulkProcessor, _ := database.DB.BulkProcessor().
		Name("bulk_post").
		Workers(2).
		BulkActions(len(datas)).
		Do(ctx)

	defer bulkProcessor.Close()

	for _, data := range datas {
		bulkProcessor.Add(
			elastic.NewBulkIndexRequest().
				Index(database.POSTINDEX).
				Id(data.Id).
				Doc(data),
		)
	}

	bulkProcessor.Flush()
	return nil
}

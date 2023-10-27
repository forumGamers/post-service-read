package documents

import (
	"context"
	"encoding/json"
	"time"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	i "github.com/forumGamers/post-service-read/interfaces"
	"github.com/olivere/elastic/v7"
)

type PostService interface {
	Insert(ctx context.Context, data PostDocument) error
	FindById(ctx context.Context, id string) (json.RawMessage, error)
	DeleteOneById(ctx context.Context, id string) error
	GetPublicContent(ctx context.Context) ([]i.PostResponse, error)
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

func (p *PostDocument) GetPublicContent(ctx context.Context) ([]i.PostResponse, error) {
	result, err := database.DB.Search().
		Index(database.POSTINDEX).
		Query(elastic.NewMatchAllQuery()).
		Size(25).Sort("createdAt", false).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var posts []PostDocument
	if result.Hits.TotalHits.Value > 0 {
		for _, hit := range result.Hits.Hits {
			var post PostDocument
			json.Unmarshal(hit.Source, &post)
			posts = append(posts, post)
		}
	}

	var postResponses []i.PostResponse
	for _, post := range posts {
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
		})
	}

	return postResponses, nil
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

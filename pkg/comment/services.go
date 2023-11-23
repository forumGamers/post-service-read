package comment

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	i "github.com/forumGamers/post-service-read/interfaces"
	"github.com/forumGamers/post-service-read/pkg/post"
	"github.com/forumGamers/post-service-read/web"
	"github.com/olivere/elastic/v7"
)

func (c *BaseDocument) Insert(ctx context.Context, data CommentDocument) error {
	return database.NewIndex(database.COMMENTINDEX).InsertOne(ctx, data)
}

func (c *BaseDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.COMMENTINDEX).DeleteOne(ctx, id)
}

func (c *BaseDocument) CountComments(ctx context.Context, posts *[]post.PostResponse, ids ...any) error {
	aggs, err := database.
		NewIndex(database.COMMENTINDEX).
		CountDocuments("postId", "comments_per_post", ids...).Do(ctx)
	if err != nil {
		return err
	}
	//hitung reply nya juga
	aggsResult, found := aggs.Aggregations.Terms("comments_per_post")
	if !found {
		return h.NotFound
	}

	for _, bucket := range aggsResult.Buckets {
		for i := 0; i < len(*posts); i++ {
			if (*posts)[i].Id == bucket.Key.(string) {
				(*posts)[i].CountComment = int(bucket.DocCount)
			}
		}
	}
	return nil
}

func (c *BaseDocument) BulkCreate(ctx context.Context, datas []CommentDocument) error {
	bulkProcessor, _ := database.DB.BulkProcessor().
		Name("bulk_comment").
		Workers(2).
		BulkActions(len(datas)).
		Do(ctx)

	defer bulkProcessor.Close()

	for _, data := range datas {
		bulkProcessor.Add(
			elastic.NewBulkIndexRequest().
				Index(database.COMMENTINDEX).
				Id(data.Id).
				Doc(data),
		)
	}

	bulkProcessor.Flush()
	return nil
}

func (c *BaseDocument) FindCommentByPostId(ctx context.Context, id string, params web.Params) ([]CommentResponse, i.TotalData, error) {
	query := database.DB.Search().
		Index(database.COMMENTINDEX).
		Query(elastic.NewMatchQuery("postId", id)).
		Sort("CreatedAt", false).
		Sort("id", false).
		Size(params.Limit).
		Timeout("30s")

	if len(params.Page) > 0 {
		var sa []any
		if timeStamp, err := strconv.ParseInt(params.Page[0], 10, 64); err == nil {
			sa = append(sa, timeStamp, params.Page[1])
			query.SearchAfter(sa...)
		}
	}

	result, err := query.
		Do(ctx)

	if err != nil {
		return nil, i.TotalData{}, err
	}

	var comments []CommentResponse
	if result.Hits.TotalHits.Value > 0 {
		for _, hit := range result.Hits.Hits {
			var comment CommentDocument
			json.Unmarshal(hit.Source, &comment)
			comments = append(comments, CommentResponse{
				Id:          comment.Id,
				UserId:      comment.UserId,
				PostId:      comment.PostId,
				Text:        h.Decryption(comment.Text),
				CreatedAt:   comment.CreatedAt,
				UpdatedAt:   comment.UpdatedAt,
				SearchAfter: hit.Sort,
			})
		}
	}

	if len(comments) < 1 {
		return comments, i.TotalData{Total: 0, Relation: "eq"}, &elastic.Error{Status: 404}
	}

	return comments, i.TotalData{Total: int(result.Hits.TotalHits.Value), Relation: result.Hits.TotalHits.Relation}, nil
}

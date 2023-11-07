package documents

import (
	"context"
	"time"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	i "github.com/forumGamers/post-service-read/interfaces"
	"github.com/olivere/elastic/v7"
)

type CommentService interface {
	Insert(ctx context.Context, data CommentDocument) error
	DeleteOneById(ctx context.Context, id string) error
	CountComments(ctx context.Context, posts *[]i.PostResponse, ids ...any) error
	BulkCreate(ctx context.Context, datas []CommentDocument) error
}

type CommentDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewComment() CommentService {
	return &CommentDocument{}
}

func (c *CommentDocument) Insert(ctx context.Context, data CommentDocument) error {
	return database.NewIndex(database.COMMENTINDEX).InsertOne(ctx, data)
}

func (c *CommentDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.COMMENTINDEX).DeleteOne(ctx, id)
}

func (c *CommentDocument) CountComments(ctx context.Context, posts *[]i.PostResponse, ids ...any) error {
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

func (c *CommentDocument) BulkCreate(ctx context.Context, datas []CommentDocument) error {
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

package documents

import (
	"context"
	"time"

	"github.com/forumGamers/post-service-read/database"
	i "github.com/forumGamers/post-service-read/interfaces"
)

type CommentService interface {
	Insert(ctx context.Context, data CommentDocument) error
	DeleteOneById(ctx context.Context, id string) error
	CountComments(ctx context.Context, posts *[]i.PostResponse, ids ...any) error
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
	aggsResult, err := database.
		NewIndex(database.COMMENTINDEX).
		CountDocuments(ctx, "postId", "comments_per_post", ids...)
	if err != nil {
		return err
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

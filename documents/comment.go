package documents

import (
	"context"
	"time"

	"github.com/forumGamers/post-service-read/database"
)

type CommentService interface {
	Insert(ctx context.Context, data CommentDocument) error
	DeleteOneById(ctx context.Context, id string) error
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

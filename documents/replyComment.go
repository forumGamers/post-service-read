package documents

import (
	"context"
	"time"

	"github.com/forumGamers/post-service-read/database"
)

type ReplyService interface {
	Insert(ctx context.Context, data ReplyCommentDocument) error
	DeleteOneById(ctx context.Context, id string) error
}

type ReplyCommentDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CommentId string `json:"commentId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *ReplyCommentDocument) Insert(ctx context.Context, data ReplyCommentDocument) error {
	return database.NewIndex(database.REPLYINDEX).InsertOne(ctx, data)
}

func (r *ReplyCommentDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.REPLYINDEX).DeleteOne(ctx, id)
}

func NewReply() ReplyService {
	return &ReplyCommentDocument{}
}

package documents

import (
	"context"
	"time"

	"github.com/forumGamers/post-service-read/database"
)

type LikeService interface {
	Insert(ctx context.Context, data LikeDocument) error
	DeleteOneById(ctx context.Context, id string) error
}

type LikeDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewLike() LikeService {
	return &LikeDocument{}
}

func (l *LikeDocument) Insert(ctx context.Context, data LikeDocument) error {
	return database.NewIndex(database.LIKEINDEX).InsertOne(ctx, data)
}

func (l *LikeDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.LIKEINDEX).DeleteOne(ctx, id)
}

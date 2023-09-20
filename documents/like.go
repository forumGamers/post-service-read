package documents

import (
	"context"
	"time"

	"github.com/forumGamers/post-service-read/database"
	i "github.com/forumGamers/post-service-read/interfaces"
)

type LikeService interface {
	Insert(ctx context.Context, data LikeDocument) error
	DeleteOneById(ctx context.Context, id string) error
	CountLike(ctx context.Context, posts *[]i.PostResponse, ids ...any) error
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

func (l *LikeDocument) CountLike(ctx context.Context, posts *[]i.PostResponse, ids ...any) error {
	aggsResult, err := database.
		NewIndex(database.COMMENTINDEX).
		CountDocuments(ctx, "postId", "likes_per_post", ids...)
	if err != nil {
		return err
	}

	for _, bucket := range aggsResult.Buckets {
		for i := 0; i < len(*posts); i++ {
			if (*posts)[i].Id == bucket.Key.(string) {
				(*posts)[i].CountLike = int(bucket.DocCount)
			}
		}
	}
	return nil
}

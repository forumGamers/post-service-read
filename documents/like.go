package documents

import (
	"context"
	"time"

	"github.com/forumGamers/post-service-read/database"
)

type LikeService interface {
	Insert(ctx context.Context, data LikeDocument) error
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
	if _, err := database.DB.
		Index().
		Index(database.LIKEINDEX).
		Id(data.Id).
		BodyJson(data).
		Do(ctx); err != nil {
		return err
	}
	return nil
}

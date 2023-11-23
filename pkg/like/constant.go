package like

import (
	"context"

	"github.com/forumGamers/post-service-read/pkg/post"
)

type LikeService interface {
	Insert(ctx context.Context, data LikeDocument) error
	DeleteOneById(ctx context.Context, id string) error
	CountLike(ctx context.Context, posts *[]post.PostResponse, userId string, ids ...any) error
	BulkCreate(ctx context.Context, datas []LikeDocument) error
}

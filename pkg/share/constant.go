package share

import (
	"context"

	"github.com/forumGamers/post-service-read/pkg/post"
)

type ShareService interface {
	Insert(ctx context.Context, data ShareDocument) error
	DeleteOneById(ctx context.Context, id string) error
	CountShares(ctx context.Context, posts *[]post.PostResponse, userId string, ids ...any) error
}

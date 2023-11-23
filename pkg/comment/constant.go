package comment

import (
	"context"

	i "github.com/forumGamers/post-service-read/interfaces"
	"github.com/forumGamers/post-service-read/pkg/post"
	"github.com/forumGamers/post-service-read/web"
)

type CommentService interface {
	Insert(ctx context.Context, data CommentDocument) error
	DeleteOneById(ctx context.Context, id string) error
	CountComments(ctx context.Context, posts *[]post.PostResponse, ids ...any) error
	BulkCreate(ctx context.Context, datas []CommentDocument) error
	FindCommentByPostId(ctx context.Context, id string, params web.Params) ([]CommentResponse, i.TotalData, error)
}

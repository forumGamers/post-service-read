package reply

import "context"

type ReplyService interface {
	Insert(ctx context.Context, data Reply) error
	DeleteOneById(ctx context.Context, id string) error
	FindCommentsReply(ctx context.Context, ids []any) ([]Reply, error)
}

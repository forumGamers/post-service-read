package post

import (
	"context"
	"encoding/json"

	i "github.com/forumGamers/post-service-read/interfaces"
	"github.com/forumGamers/post-service-read/web"
)

type PostService interface {
	Insert(ctx context.Context, data PostDocument) error
	FindById(ctx context.Context, id string) (json.RawMessage, error)
	DeleteOneById(ctx context.Context, id string) error
	GetPublicContent(ctx context.Context, query web.PostParams) ([]PostResponse, i.TotalData, error)
	BulkCreate(ctx context.Context, datas []PostDocument) error
	FindByUserId(ctx context.Context, id string) ([]PostResponse, i.TotalData, error)
}

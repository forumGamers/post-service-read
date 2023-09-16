package documents

import (
	"context"
	"time"

	"github.com/forumGamers/post-service-read/database"
)

type ShareService interface {
	Insert(ctx context.Context, data ShareDocument) error
	DeleteOneById(ctx context.Context, id string) error
}

type ShareDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	Text      string `json:"text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *ShareDocument) Insert(ctx context.Context, data ShareDocument) error {
	return database.NewIndex(database.SHAREINDEX).InsertOne(ctx, data)
}

func (s *ShareDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.SHAREINDEX).DeleteOne(ctx, id)
}

func NewShare() ShareService {
	return &ShareDocument{}
}

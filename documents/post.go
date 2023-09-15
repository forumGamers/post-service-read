package documents

import (
	"context"
	"encoding/json"
	"time"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/olivere/elastic/v7"
)

type PostService interface {
	Insert(ctx context.Context, data PostDocument) error
	FindById(ctx context.Context, id string) (json.RawMessage, error)
	DeleteOneById(ctx context.Context, id string) error
}

type Media struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PostDocument struct {
	Id           string `json:"id"`
	UserId       string `json:"userId"`
	Text         string `json:"text" bson:"text"`
	Media        Media
	AllowComment bool `json:"allowComment" default:"true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tags         []string `json:"tags"`
	Privacy      string   `json:"privacy" default:"Public"`
}

func (p *PostDocument) Insert(ctx context.Context, data PostDocument) error {
	if _, err := database.DB.
		Index().
		Index(database.POSTINDEX).
		Id(data.Id).
		BodyJson(data).
		Do(ctx); err != nil {
		return err
	}
	return nil
}

func (p *PostDocument) DeleteOneById(ctx context.Context, id string) error {
	if _, err := database.DB.Delete().
		Index(database.POSTINDEX).
		Id(id).
		Do(ctx); err != nil {
		if elastic.IsNotFound(err) {
			return h.NotFound
		}
		return err

	}
	return nil
}

func NewPost() PostService {
	return &PostDocument{}
}

func (p *PostDocument) FindById(ctx context.Context, id string) (json.RawMessage, error) {
	get, err := database.DB.Get().Index(database.POSTINDEX).Id(id).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, h.NotFound
		}

		return nil, err
	}

	return get.Source, nil
}
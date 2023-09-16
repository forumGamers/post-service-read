package database

import (
	"context"
	"reflect"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/olivere/elastic/v7"
)

type Index struct {
	Name string
}

type Operations interface {
	InsertOne(ctx context.Context, data any) error
	DeleteOne(ctx context.Context, id string) error
}

func NewIndex(name string) Operations {
	return &Index{name}
}

func (i *Index) InsertOne(ctx context.Context, data any) error {
	if _, err := DB.Index().
		Index(i.Name).
		Id(reflect.ValueOf(data).FieldByName("Id").String()).
		BodyJson(data).
		Do(ctx); err != nil {
		return err
	}
	return nil
}

func (i *Index) DeleteOne(ctx context.Context, id string) error {
	if _, err := DB.Delete().Index(i.Name).Id(id).Do(ctx); err != nil {
		if elastic.IsNotFound(err) {
			return h.NotFound
		}
		return err
	}
	return nil
}

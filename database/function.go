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
	CountDocuments(termsField, aggregationName string, ids ...any) *elastic.SearchService
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

func (i *Index) CountDocuments(termsField, aggregationName string, ids ...any) *elastic.SearchService {
	return DB.Search().
		Index(i.Name).
		Size(0).
		Query(elastic.NewTermsQuery(termsField, ids...)).
		Aggregation(aggregationName, elastic.NewTermsAggregation().Field(termsField))
}

package post

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	i "github.com/forumGamers/post-service-read/interfaces"
	"github.com/forumGamers/post-service-read/web"
	"github.com/olivere/elastic/v7"
)

func (p *BaseDocument) Insert(ctx context.Context, data PostDocument) error {
	return database.NewIndex(database.POSTINDEX).InsertOne(ctx, data)
}

func (p *BaseDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.POSTINDEX).DeleteOne(ctx, id)
}

func (p *BaseDocument) FindById(ctx context.Context, id string) (json.RawMessage, error) {
	get, err := p.DB.Get().Index(database.POSTINDEX).Id(id).Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, h.NotFound
		}

		return nil, err
	}

	return get.Source, nil
}

func (p *BaseDocument) GetPublicContent(ctx context.Context, query web.PostParams) ([]PostResponse, i.TotalData, error) {
	search := p.DB.Search().
		Index(database.POSTINDEX).
		Size(query.Limit).
		Sort("CreatedAt", false).
		Sort("id", false)

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewRangeQuery("CreatedAt").Gte("now-30d/d").Lte("now/d"))

	if len(query.UserIds) > 0 {
		var ids []any
		for _, val := range query.UserIds {
			ids = append(ids, val)
		}
		boolQuery.Must(elastic.NewTermsQuery("userId", ids...))
	}
	search.Query(boolQuery)

	if len(query.Page) > 0 || query.Page != nil {
		var sa []any
		if timeStamp, err := strconv.ParseInt(query.Page[0], 10, 64); err == nil {
			sa = append(sa, timeStamp, query.Page[1])
			search.SearchAfter(sa...)
		}
	}
	search.Timeout("30s")

	result, err := search.Do(context.Background())
	if err != nil {
		return nil, struct {
			Total    int
			Relation string
		}{}, err
	}

	var postResponses []PostResponse
	if result.Hits.TotalHits.Value > 0 {
		ParseToPostResponse(result.Hits.Hits, &postResponses)
	}

	if len(postResponses) < 1 {
		return postResponses, i.TotalData{Total: 0, Relation: "eq"}, &elastic.Error{Status: 404}
	}

	return postResponses, i.TotalData{Total: int(result.Hits.TotalHits.Value), Relation: result.Hits.TotalHits.Relation}, nil
}

func (p *BaseDocument) BulkCreate(ctx context.Context, datas []PostDocument) error {
	bulkProcessor, _ := p.DB.BulkProcessor().
		Name("bulk_post").
		Workers(2).
		BulkActions(len(datas)).
		Do(ctx)

	defer bulkProcessor.Close()

	for _, data := range datas {
		bulkProcessor.Add(
			elastic.NewBulkIndexRequest().
				Index(database.POSTINDEX).
				Id(data.Id).
				Doc(data),
		)
	}

	bulkProcessor.Flush()
	return nil
}

func (p *BaseDocument) FindByUserId(ctx context.Context, id string) ([]PostResponse, i.TotalData, error) {
	result, err := p.DB.Search().
		Index(database.POSTINDEX).
		Query(
			elastic.NewBoolQuery().
				Must(elastic.NewMatchQuery("userId", id)),
		).
		Size(10).
		Sort("CreatedAt", false).
		Sort("id", false).
		Do(ctx)
	if err != nil {
		return nil, struct {
			Total    int
			Relation string
		}{}, err
	}

	var postResponses []PostResponse
	if result.Hits.TotalHits.Value > 0 {
		ParseToPostResponse(result.Hits.Hits, &postResponses)
	}
	if len(postResponses) < 1 {
		return postResponses, i.TotalData{Total: 0, Relation: "eq"}, &elastic.Error{Status: 404}
	}

	return postResponses, i.TotalData{Total: int(result.Hits.TotalHits.Value), Relation: result.Hits.TotalHits.Relation}, nil
}

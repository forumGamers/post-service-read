package reply

import (
	"context"

	"github.com/forumGamers/post-service-read/database"
	"github.com/olivere/elastic/v7"
)

func (r *BaseDocument) Insert(ctx context.Context, data Reply) error {
	return database.NewIndex(database.REPLYINDEX).InsertOne(ctx, data)
}

func (r *BaseDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.REPLYINDEX).DeleteOne(ctx, id)
}

func (r *BaseDocument) FindCommentsReply(ctx context.Context, ids []any) ([]Reply, error) {
	result, err := database.DB.Search().
		Index(database.REPLYINDEX).
		Query(elastic.NewBoolQuery().Must(elastic.NewTermsQuery("commentId", ids...))).
		Timeout("30s").
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var replies []Reply
	if result.Hits.TotalHits.Value > 0 {
		ParseToReply(result.Hits.Hits, &replies)
	}
	return replies, nil
}

package like

import (
	"context"
	"sync"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/post"
	"github.com/olivere/elastic/v7"
)

func (l *BaseDocument) Insert(ctx context.Context, data LikeDocument) error {
	return database.NewIndex(database.LIKEINDEX).InsertOne(ctx, data)
}

func (l *BaseDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.LIKEINDEX).DeleteOne(ctx, id)
}

func (l *BaseDocument) CountLike(ctx context.Context, posts *[]post.PostResponse, userId string, ids ...any) error {
	query := database.
		NewIndex(database.LIKEINDEX).
		CountDocuments("postId", "likes_per_post", ids...)

	if userId != "" {
		query = query.Aggregation(
			"liked", elastic.NewFilterAggregation().
				Filter(elastic.NewTermQuery("userId", userId)).
				SubAggregation("postId", elastic.NewTermsAggregation().Field("postId")),
		)
	}

	result, err := query.Do(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	if userId != "" {
		filterLiked, found := result.Aggregations.Filter("liked")
		if !found {
			return h.NotFound
		}

		aggsLiked, found := filterLiked.Terms("postId")
		if !found {
			return h.NotFound
		}

		wg.Add(1)
		go func(postsData *[]post.PostResponse) {
			defer wg.Done()
			for _, bucket := range aggsLiked.Buckets {
				for i := 0; i < len(*postsData); i++ {
					if (*postsData)[i].Id == bucket.Key.(string) {
						(*postsData)[i].IsLiked = true
					}
				}
			}
		}(posts)
	}

	aggsLikePerPost, found := result.Aggregations.Terms("likes_per_post")
	if !found {
		return h.NotFound
	}

	wg.Add(1)
	go func(postsData *[]post.PostResponse) {
		defer wg.Done()
		for _, bucket := range aggsLikePerPost.Buckets {
			for i := 0; i < len(*postsData); i++ {
				if (*postsData)[i].Id == bucket.Key.(string) {
					(*postsData)[i].CountLike = int(bucket.DocCount)
				}
			}
		}
	}(posts)

	wg.Wait()
	return nil
}

func (l *BaseDocument) BulkCreate(ctx context.Context, datas []LikeDocument) error {
	bulkProcessor, _ := database.DB.BulkProcessor().
		Name("bulk_like").
		Workers(2).
		BulkActions(len(datas)).
		Do(ctx)

	defer bulkProcessor.Close()

	for _, data := range datas {
		bulkProcessor.Add(
			elastic.NewBulkIndexRequest().
				Index(database.LIKEINDEX).
				Id(data.Id).
				Doc(data),
		)
	}

	bulkProcessor.Flush()
	return nil
}

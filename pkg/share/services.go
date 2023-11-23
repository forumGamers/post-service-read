package share

import (
	"context"
	"sync"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/post"
	"github.com/olivere/elastic/v7"
)

func (s *BaseDocument) Insert(ctx context.Context, data ShareDocument) error {
	return database.NewIndex(database.SHAREINDEX).InsertOne(ctx, data)
}

func (s *BaseDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.SHAREINDEX).DeleteOne(ctx, id)
}

func (s *BaseDocument) CountShares(ctx context.Context, posts *[]post.PostResponse, userId string, ids ...any) error {
	query := database.
		NewIndex(database.SHAREINDEX).
		CountDocuments("postId", "shares_per_post", ids...)

	if userId != "" {
		query = query.Aggregation(
			"shared", elastic.NewFilterAggregation().
				Filter(elastic.NewTermQuery("userId", userId)).
				SubAggregation("postId", elastic.NewTermsAggregation().Field("postId")),
		)
	}
	// 6542168defedc06c06df5745,6542168defedc06c06df5747,6542168defedc06c06df5748,6542168defedc06c06df5751,6542168defedc06c06df5754,6542168defedc06c06df5757,6542168defedc06c06df575a
	result, err := query.Do(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	if userId != "" {
		filterShared, found := result.Aggregations.Filter("shared")
		if !found {
			return h.NotFound
		}

		aggsShared, found := filterShared.Terms("postId")
		if !found {
			return h.NotFound
		}

		wg.Add(1)
		go func(postsData *[]post.PostResponse) {
			defer wg.Done()
			for _, bucket := range aggsShared.Buckets {
				for i := 0; i < len(*postsData); i++ {
					if (*postsData)[i].Id == bucket.Key.(string) {
						(*postsData)[i].IsShared = true
					}
				}
			}
		}(posts)
	}

	aggsSharesPerPost, found := result.Aggregations.Terms("shares_per_post")
	if !found {
		return h.NotFound
	}

	wg.Add(1)
	go func(postsData *[]post.PostResponse) {
		defer wg.Done()
		for _, bucket := range aggsSharesPerPost.Buckets {
			for i := 0; i < len(*postsData); i++ {
				if (*postsData)[i].Id == bucket.Key.(string) {
					(*postsData)[i].CountShare = int(bucket.DocCount)
				}
			}
		}
	}(posts)

	wg.Wait()
	return nil
}

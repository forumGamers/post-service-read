package documents

import (
	"context"
	"sync"
	"time"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	i "github.com/forumGamers/post-service-read/interfaces"
	"github.com/olivere/elastic/v7"
)

type LikeService interface {
	Insert(ctx context.Context, data LikeDocument) error
	DeleteOneById(ctx context.Context, id string) error
	CountLike(ctx context.Context, posts *[]i.PostResponse, userId string, ids ...any) error
}

type LikeDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewLike() LikeService {
	return &LikeDocument{}
}

func (l *LikeDocument) Insert(ctx context.Context, data LikeDocument) error {
	return database.NewIndex(database.LIKEINDEX).InsertOne(ctx, data)
}

func (l *LikeDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.LIKEINDEX).DeleteOne(ctx, id)
}

func (l *LikeDocument) CountLike(ctx context.Context, posts *[]i.PostResponse, userId string, ids ...any) error {
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
		go func(postsData *[]i.PostResponse) {
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
	go func(postsData *[]i.PostResponse) {
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

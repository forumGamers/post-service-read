package documents

import (
	"context"
	"encoding/json"
	"time"

	"github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	i "github.com/forumGamers/post-service-read/interfaces"
	"github.com/olivere/elastic/v7"
)

type ReplyService interface {
	Insert(ctx context.Context, data ReplyCommentDocument) error
	DeleteOneById(ctx context.Context, id string) error
	FindCommentsReply(ctx context.Context, ids []any) ([]i.Reply, error)
}

type ReplyCommentDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CommentId string `json:"commentId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *ReplyCommentDocument) Insert(ctx context.Context, data ReplyCommentDocument) error {
	return database.NewIndex(database.REPLYINDEX).InsertOne(ctx, data)
}

func (r *ReplyCommentDocument) DeleteOneById(ctx context.Context, id string) error {
	return database.NewIndex(database.REPLYINDEX).DeleteOne(ctx, id)
}

func (r *ReplyCommentDocument) FindCommentsReply(ctx context.Context, ids []any) ([]i.Reply, error) {
	result, err := database.DB.Search().
		Index(database.REPLYINDEX).
		Query(elastic.NewBoolQuery().Must(elastic.NewTermsQuery("commentId", ids...))).
		Timeout("30s").
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var replies []i.Reply
	if result.Hits.TotalHits.Value > 0 {
		for _, hit := range result.Hits.Hits {
			var reply ReplyCommentDocument
			json.Unmarshal(hit.Source, &reply)
			replies = append(replies, i.Reply{
				Text:      h.Decryption(reply.Text),
				Id:        reply.Id,
				UserId:    reply.UserId,
				CommentId: reply.CommentId,
				CreatedAt: reply.CreatedAt,
				UpdatedAt: reply.UpdatedAt,
			})
		}
	}
	return replies, nil
}

func NewReply() ReplyService {
	return &ReplyCommentDocument{}
}

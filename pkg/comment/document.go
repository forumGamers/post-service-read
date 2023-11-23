package comment

import (
	"time"

	"github.com/forumGamers/post-service-read/database"
	"github.com/forumGamers/post-service-read/pkg/reply"
	"github.com/olivere/elastic/v7"
)

type CommentDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentResponse struct {
	Id          string `json:"_id"`
	UserId      string `json:"userId"`
	Text        string `json:"text"`
	PostId      string `json:"postId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SearchAfter []any         `json:"searchAfter"`
	Reply       []reply.Reply `json:"Reply"`
}

type BaseDocument struct {
	DB *elastic.Client
}

func NewComment() CommentService {
	return &BaseDocument{
		DB: database.DB,
	}
}

package reply

import (
	"time"

	"github.com/forumGamers/post-service-read/database"
	"github.com/olivere/elastic/v7"
)

type Reply struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CommentId string `json:"commentId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseDocument struct {
	DB *elastic.Client
}

func NewReply() ReplyService {
	return &BaseDocument{
		DB: database.DB,
	}
}

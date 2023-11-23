package share

import (
	"time"

	"github.com/forumGamers/post-service-read/database"
	"github.com/olivere/elastic/v7"
)

type ShareDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	Text      string `json:"text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseDocument struct {
	DB *elastic.Client
}

func NewShare() ShareService {
	return &BaseDocument{
		DB: database.DB,
	}
}

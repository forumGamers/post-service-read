package like

import (
	"time"

	"github.com/forumGamers/post-service-read/database"
	"github.com/olivere/elastic/v7"
)

type LikeDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BaseDocument struct {
	DB *elastic.Client
}

func NewLike() LikeService {
	return &BaseDocument{
		DB: database.DB,
	}
}

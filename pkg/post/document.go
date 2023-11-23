package post

import (
	"time"

	"github.com/forumGamers/post-service-read/database"
	"github.com/olivere/elastic/v7"
)

type PostDocument struct {
	Id           string `json:"id"`
	UserId       string `json:"userId"`
	Text         string `json:"text" bson:"text"`
	Media        Media
	AllowComment bool `json:"allowComment" default:"true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tags         []string `json:"tags"`
	Privacy      string   `json:"privacy" default:"Public"`
}

type Media struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PostResponse struct {
	Id           string `json:"_id"`
	UserId       string `json:"userId"`
	Text         string `json:"text"`
	Media        Media
	AllowComment bool `json:"allowComment"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CountLike    int
	CountComment int
	CountShare   int
	IsLiked      bool     `json:"isLiked"`
	IsShared     bool     `json:"isShared"`
	Tags         []string `json:"tags"`
	Privacy      string   `json:"privacy"`
	SearchAfter  []any    `json:"searchAfter"`
}

type BaseDocument struct {
	DB *elastic.Client
}

func NewPost() PostService {
	return &BaseDocument{
		DB: database.DB,
	}
}

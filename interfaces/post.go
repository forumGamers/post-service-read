package interfaces

import (
	"time"
)

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

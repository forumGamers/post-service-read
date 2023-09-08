package models

import (
	"time"
)

type Media struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	Id   string `json:"id"`
}

type Post struct {
	Id           string `json:"_id"`
	UserId       string `json:"userId"`
	Text         string `json:"text" bson:"text"`
	Media        Media
	AllowComment bool `json:"allowComment" default:"true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tags         []string `json:"tags"`
	Privacy      string   `json:"privacy" default:"Public"`
}

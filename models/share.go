package models

import (
	"time"
)

type Share struct {
	Id        string `json:"_id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	Text      string `json:"text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

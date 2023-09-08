package models

import (
	"time"
)

type Comment struct {
	Id        string `json:"_id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

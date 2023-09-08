package models

import (
	"time"
)

type ReplyComment struct {
	Id        string `json:"_id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CommentId string `json:"commentId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

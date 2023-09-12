package documents

import (
	"time"
)

type CommentDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

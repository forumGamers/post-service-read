package documents

import (
	"time"
)

type ReplyCommentDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CommentId string `json:"commentId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

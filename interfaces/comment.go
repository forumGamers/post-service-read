package interfaces

import (
	"time"
)

type CommentResponse struct {
	Id          string `json:"_id"`
	UserId      string `json:"userId"`
	Text        string `json:"text"`
	PostId      string `json:"postId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SearchAfter []any   `json:"searchAfter"`
	Reply       []Reply `json:"reply"`
}

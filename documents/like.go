package documents

import (
	"time"
)

type LikeDocument struct {
	Id        string `json:"_id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

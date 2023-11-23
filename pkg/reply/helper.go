package reply

import (
	"encoding/json"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/olivere/elastic/v7"
)

func ParseToReply(source []*elastic.SearchHit, result *[]Reply) {
	for _, hit := range source {
		var reply Reply
		json.Unmarshal(hit.Source, &reply)
		*result = append(*result, Reply{
			Text:      h.Decryption(reply.Text),
			Id:        reply.Id,
			UserId:    reply.UserId,
			CommentId: reply.CommentId,
			CreatedAt: reply.CreatedAt,
			UpdatedAt: reply.UpdatedAt,
		})
	}
}

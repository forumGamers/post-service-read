package post

import (
	"encoding/json"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/olivere/elastic/v7"
)

func ParseToPostResponse(source []*elastic.SearchHit, result *[]PostResponse) {
	for _, hit := range source {
		var post PostDocument
		json.Unmarshal(hit.Source, &post)
		*result = append(*result, PostResponse{
			Id:           post.Id,
			UserId:       post.UserId,
			Text:         h.Decryption(post.Text),
			Media:        Media(post.Media),
			AllowComment: post.AllowComment,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			Tags:         post.Tags,
			Privacy:      post.Privacy,
			SearchAfter:  hit.Sort,
		})
	}
}

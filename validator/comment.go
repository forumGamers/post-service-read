package validator

import "github.com/forumGamers/post-service-read/web"

func DefaultLimit(query *web.Params) {
	if query.Limit == 0 {
		query.Limit = 10
	}
}

func ValidateCommentParams(query *web.Params) {
	DefaultLimit(query)
	ParsePageQuery(&query.Page)
}

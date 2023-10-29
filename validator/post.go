package validator

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/forumGamers/post-service-read/web"
)

const (
	regexID = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`
)

func ValidatePaginations(query *web.PostParams) {
	if query.Limit == 0 {
		query.Limit = 20
	}
}

func ValidateUserId(query *web.PostParams) {
	var ids []string
	if query.UserIds[0] == "" {
		return
	}

	if !strings.Contains(query.UserIds[0], ",") {
		return
	}

	for _, val := range strings.Split(query.UserIds[0], ",") {
		if !contains(ids, val) &&
			regexp.MustCompile(regexID).MatchString(val) &&
			val != "" {
			ids = append(ids, val)
		}
	}
	query.UserIds = ids
}

func ValidatePostQuery(query *web.PostParams) {
	ValidatePaginations(query)
	ValidateUserId(query)
	ParsePageQuery(query)
}

func contains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func ParsePageQuery(query *web.PostParams) {
	if query.Page == nil {
		return
	}

	partStr, ok := query.Page[0].(string)
	if !ok {
		query.Page = nil
		return
	}

	parts := strings.Split(partStr, ",")
	if len(parts) < 2 {
		query.Page = nil
		return
	}

	timeStamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		query.Page = nil
		return
	}

	query.Page = []any{timeStamp, parts[1]}
}

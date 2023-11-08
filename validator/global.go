package validator

import (
	"strings"
)

func ParsePageQuery(query *[]string) {
	if query == nil || len(*query) < 1 {
		return
	}
	partStr := (*query)[0]

	parts := strings.Split(partStr, ",")
	if len(parts) < 2 {
		query = nil
		return
	}

	*query = parts
}

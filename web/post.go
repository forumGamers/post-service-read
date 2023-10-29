package web

type PostParams struct {
	UserIds    []string `form:"userIds"`
	Page       []string `form:"page"`
	Limit      int      `form:"limit"`
	Sort       string   `form:"sort"`
	Preference string   `form:"preference"`
}

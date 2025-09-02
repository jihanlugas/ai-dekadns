package request

type Pagination struct {
	Page  string `json:"page" form:"page" query:"page" default:"1"`
	Limit string `json:"limit" form:"limit" query:"limit" default:"10"`
}

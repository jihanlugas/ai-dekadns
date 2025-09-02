package request

type PageZone struct {
	Pagination
	ProjectId string `json:"projectId" query:"projectId" form:"projectId" validate:"required"`
}
type CreateZone struct {
	Name      string `json:"name" validate:"domain"`
	ProjectId string `json:"projectId" validate:"required"`
}

package request

type CreateDns struct {
	Name      string `json:"name" validate:"required"`
	ProjectId string `json:"projectId" validate:"required"`
}

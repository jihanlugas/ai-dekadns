package request

type CreateDns struct {
	Name      string `json:"name" validate:"domain"`
	ProjectId string `json:"projectId" validate:"required"`
}

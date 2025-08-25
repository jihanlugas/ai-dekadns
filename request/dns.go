package request

type CreateDns struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	OrganizationId string `json:"organizationId"`
	ProjectId      string `json:"projectId"`
	IsCustomNs     bool   `json:"isCustomNs"`
	IsDnssec       string `json:"isDnssec"`
}

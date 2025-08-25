package request

type RequestAuditlog struct {
	UserID         string `json:"user_id"`
	ServiceName    string `json:"service_name"`
	ServiceType    string `json:"service_type"`
	MethodName     string `json:"method_name"`
	Action         string `json:"action"`
	Level          string `json:"level"`
	ProjectID      string `json:"project_id"`
	OrganizationID string `json:"organization_id"`
	IPAddress      string `json:"ipAddress"`
}

package request

type CreateRecord struct {
	ZoneID  string `json:"zoneId" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Content string `json:"content" validate:"required"`
	TTL     uint32 `json:"ttl" validate:"required"`
}

type DeleteRecord struct {
	ZoneID  string `json:"zoneId" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Type    string `json:"type" validate:"required"`
	Content string `json:"content" validate:"required"`
	TTL     uint32 `json:"ttl" validate:"required"`
}

type UpdateRecord struct {
	ZoneID     string `json:"zoneId" validate:"required"`
	OldName    string `json:"oldName" validate:"required"`
	OldType    string `json:"oldType" validate:"required"`
	OldContent string `json:"oldContent" validate:"required"`
	OldTTL     uint32 `json:"oldTtl" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Type       string `json:"type" validate:"required"`
	Content    string `json:"content" validate:"required"`
	TTL        uint32 `json:"ttl" validate:"required"`
}

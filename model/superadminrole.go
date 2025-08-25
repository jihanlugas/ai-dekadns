package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SuperAdminRoles struct {
	Id          uuid.UUID        `gorm:"type:uuid;primary_key;" json:"id"`
	Name        string           `gorm:"varchar(100)" json:"name"`
	Default     bool             `gorm:"bool;" json:"default"`
	Desc        *string          `gorm:"text" json:"desc"`
	Usages      *int             `gorm:"-" json:"usages,omitempty"`
	CreatedAt   time.Time        `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	Privileges  *AdminPrivileges `gorm:"type:text;column:privileges;" json:"privileges,omitempty"`
	SpesificOrg string           `gorm:"type:text;column:spesific_org;" json:"spesific_org,omitempty"`
}

type AdminPrivilege struct {
	Name  string   `json:"name"`
	Role  string   `json:"role"`
	Types []string `json:"types,omitempty"`
}

type AdminPrivileges []AdminPrivilege

func (d *AdminPrivileges) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &d)
	case string:
		return json.Unmarshal([]byte(v), &d)
	}
	return errors.New("type assertion failed")
}

func (d AdminPrivileges) Value() (driver.Value, error) {
	val, err := json.Marshal(d)
	return val, err
}

func (d *AdminPrivileges) GetTypes(privilegeName string) (orgTypes []string) {
	for _, privilege := range *d {
		if privilege.Name == privilegeName {
			orgTypes = privilege.Types
			break
		}
	}

	return
}

func (s *SuperAdminRoles) CheckForPrivilege(privilegeName string, requireEditorAccess bool) bool {
	if s.Privileges != nil {
		for _, privilege := range *s.Privileges {
			if strings.EqualFold(privilege.Name, privilegeName) {

				if requireEditorAccess && !strings.EqualFold(privilege.Role, "Editor") {
					return false
				}

				return true
			}
		}
	}

	return false
}

func (SuperAdminRoles) TableName() string {
	return "superadminroles"
}

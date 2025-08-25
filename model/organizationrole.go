package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OrganizationRole struct {
	ID             int64          `json:"id" gorm:"primary_key AUTO_INCREMENT"`
	OrganizationID string         `json:"organization_id,omitempty"`
	Name           string         `gorm:"not null" json:"name" `
	Description    string         `json:"description"`
	Privilages     Privilages     `json:"privilages" gorm:"type:json"`
	IsDefault      bool           `json:"is_default"`
	CreatedAt      time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// CheckForPrivilege checks if the user has the required privileges.
// If privilegeName not found (case-sensitive), the function returns false.
// If requireEditorAccess is true, the user must have the editor privilege.
func (o *OrganizationRole) CheckForPrivilege(privilegeName string, requireEditorAccess bool) bool {
	for _, privilege := range o.Privilages {
		// Find for privilege.
		if privilege.Name == privilegeName {
			// If disabled then return false.
			if privilege.Disable {
				return false
			}

			// Check if requireEditorAccess.
			if requireEditorAccess && !privilege.Editor {
				return false
			}

			// Else is viewer. No need to check.
			return true
		}
	}

	// Privilege not found.
	return false
}

type Privilages []struct {
	Disable    bool   `json:"disable"`
	Name       string `gorm:"not null" json:"name" `
	Editor     bool   `json:"editor"`
	Viewer     bool   `json:"viewer"`
	Suspend    bool   `json:"suspend"`
	Terminated bool   `json:"terminated"`
}

func (d *Privilages) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &d)
	case string:
		return json.Unmarshal([]byte(v), &d)
	}
	return errors.New("type assertion failed")
}

func (d Privilages) Value() (driver.Value, error) {
	val, err := json.Marshal(d)
	return string(val), err
}

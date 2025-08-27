package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Zone struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	OrganizationId string         `gorm:"not null" json:"organizationId"`
	ProjectId      string         `gorm:"not null" json:"projectId"`
	Name           string         `gorm:"not null" json:"name"`
	Status         string         `gorm:"not null" json:"status"`
	IsCustomNs     bool           `gorm:"not null" json:"isCustomNs"`
	IsDnssec       string         `gorm:"not null" json:"isDnssec"`
	CreatedBy      string         `gorm:"not null" json:"createdBy"`
	UpdatedBy      string         `gorm:"not null" json:"updatedBy"`
	CreatedAt      time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"null" json:"-" `
}

func (m *Zone) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()

	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}

	return nil
}

func (m *Zone) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return
}

func (Zone) TableName() string {
	return "dekadns.zones"
}

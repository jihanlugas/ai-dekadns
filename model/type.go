package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Type struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	DeletedAt gorm.DeletedAt `gorm:"null" json:"-"`
}

func (m *Type) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	return nil
}

func (m *Type) BeforeUpdate(tx *gorm.DB) (err error) {
	return
}

func (Type) TableName() string {
	return "dekadns.types"
}

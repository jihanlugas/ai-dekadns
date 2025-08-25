package organization

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetById(id string) (org model.Organization, err error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) GetById(id string) (org model.Organization, err error) {
	conn := r.db
	err = conn.Where("id = ?", id).First(&org).Error
	return org, err
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

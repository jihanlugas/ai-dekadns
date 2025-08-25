package project

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetById(id string) (proj model.Project, err error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) GetById(id string) (proj model.Project, err error) {
	conn := r.db
	err = conn.Where("id = ?", id).First(&proj).Error
	return proj, err
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

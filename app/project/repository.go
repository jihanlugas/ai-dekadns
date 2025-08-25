package project

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetById(id string) (proj model.Project, err error)
}

type repository struct {
	coreDb *gorm.DB
	sslDb  *gorm.DB
}

func (r repository) GetById(id string) (proj model.Project, err error) {
	conn := r.coreDb
	err = conn.Where("id = ?", id).First(&proj).Error
	return proj, err
}

func NewRepository(coreDb *gorm.DB, sslDb *gorm.DB) Repository {
	return &repository{
		coreDb: coreDb,
		sslDb:  sslDb,
	}
}

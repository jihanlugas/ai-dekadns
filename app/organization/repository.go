package organization

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetById(id string) (org model.Organization, err error)
}

type repository struct {
	coreDb *gorm.DB
	sslDb  *gorm.DB
}

func (r repository) GetById(id string) (org model.Organization, err error) {
	conn := r.coreDb
	err = conn.Where("id = ?", id).First(&org).Error
	return org, err
}

func NewRepository(coreDb *gorm.DB, sslDb *gorm.DB) Repository {
	return &repository{
		coreDb: coreDb,
		sslDb:  sslDb,
	}
}

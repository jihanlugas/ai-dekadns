package user

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetById(id string, preloads ...string) (ust model.User, err error)
}

type repository struct {
	coreDb *gorm.DB
	sslDb  *gorm.DB
}

func (r repository) GetById(id string, preloads ...string) (ust model.User, err error) {
	conn := r.coreDb

	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	err = conn.Where("id = ?", id).First(&ust).Error
	return ust, err
}

func NewRepository(coreDb *gorm.DB, sslDb *gorm.DB) Repository {
	return &repository{
		coreDb: coreDb,
		sslDb:  sslDb,
	}
}

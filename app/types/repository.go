package types

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetById(id string) (data model.Type, err error)
	Create(data model.Type) (err error)
	Update(data model.Type) (err error)
	Delete(id string) (err error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) GetById(id string) (data model.Type, err error) {
	return data, r.db.First(&data, "id = ?", id).Error
}

func (r repository) Create(data model.Type) (err error) {
	return r.db.Create(&data).Error
}

func (r repository) Update(data model.Type) (err error) {
	return r.db.Save(&data).Error
}

func (r repository) Delete(id string) (err error) {
	return r.db.Delete(&model.Type{}, "id = ?", id).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

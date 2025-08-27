package record

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetById(id string) (data model.Record, err error)
	Create(data model.Record) (err error)
	Update(data model.Record) (err error)
	Delete(id string) (err error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) GetById(id string) (data model.Record, err error) {
	return data, r.db.First(&data, "id = ?", id).Error
}

func (r repository) Create(data model.Record) (err error) {
	return r.db.Create(&data).Error
}

func (r repository) Update(data model.Record) (err error) {
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

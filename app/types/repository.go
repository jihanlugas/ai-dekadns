package types

import (
	"ai-dekadns/helper"
	"ai-dekadns/model"
	"ai-dekadns/request"
	"math"

	"gorm.io/gorm"
)

type Repository interface {
	Page(req request.PageType, pageReq model.Pagination) (*model.Pagination, error)
	GetById(id string) (data model.Type, err error)
	Create(data model.Type) (err error)
	Update(data model.Type) (err error)
	Delete(id string) (err error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) Page(req request.PageType, pageReq model.Pagination) (*model.Pagination, error) {
	var err error
	conn := r.db
	var types []*model.Type
	var totalRows int64
	//var resp = make([]model.Type, 0)

	conn = conn.Model(types)

	err = conn.Scopes(helper.Paginate(&pageReq, conn)).Find(&types).Error
	if err != nil {
		return nil, err
	}

	err = conn.Count(&totalRows).Error
	if err != nil {
		return nil, err
	}

	pageReq.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pageReq.GetLimit())))
	pageReq.TotalPages = totalPages
	pageReq.Rows = types

	return &pageReq, nil
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

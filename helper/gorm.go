package helper

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

func Paginate(pagination *model.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

package superadminrole

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetBySuperadminId(id string) (rolePrivilege []model.SuperAdminRoles, err error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) GetBySuperadminId(id string) (rolePrivilege []model.SuperAdminRoles, err error) {
	conn := r.db
	return rolePrivilege, conn.Joins("JOIN super_admin_role_mapping srm ON srm.id_role = super_admin_roles.id").
		Where("srm.id_admin = ?", id).Find(&rolePrivilege).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

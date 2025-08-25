package superadminrole

import (
	"ai-dekadns/model"

	"gorm.io/gorm"
)

type Repository interface {
	GetBySuperadminId(id string) (rolePrivilege []model.SuperAdminRoles, err error)
}

type repository struct {
	coreDb *gorm.DB
	sslDb  *gorm.DB
}

func (r repository) GetBySuperadminId(id string) (rolePrivilege []model.SuperAdminRoles, err error) {
	conn := r.coreDb
	return rolePrivilege, conn.Joins("JOIN super_admin_role_mapping srm ON srm.id_role = super_admin_roles.id").
		Where("srm.id_admin = ?", id).Find(&rolePrivilege).Error
}

func NewRepository(coreDb *gorm.DB, sslDb *gorm.DB) Repository {
	return &repository{
		coreDb: coreDb,
		sslDb:  sslDb,
	}
}

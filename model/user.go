package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                 string     `gorm:"primaryKey;type:uuid" json:"id" validate:"xss_clean" conform:"trim"`
	OrganizationID     string     `json:"organization_id,omitempty" validate:"xss_clean" conform:"trim"`
	Type               string     `gorm:"null; type:char(1)" json:"type" validate:"xss_clean" conform:"trim"`
	Fullname           string     `gorm:"not null" json:"fullname" validate:"specialchar" conform:"trim"`
	Firstname          string     `json:"firstname" validate:"specialchar" conform:"trim"`
	Lastname           string     `json:"lastname" validate:"specialchar" conform:"trim"`
	Email              string     `gorm:"not null,unique" json:"email,omitempty" validate:"email" conform:"trim"`
	Password           string     `gorm:"not null" json:"password,omitempty" validate:""`
	PhoneNumber        string     `json:"phone_number" validate:"xss_clean" conform:"trim"`
	NpwpPrivacy        string     `json:"npwp_privacy" validate:"xss_clean" conform:"trim"`
	RoleID             int64      `json:"role_id,omitempty" validate:""`
	Role               Role       `json:"role" validate:""`
	SuperAdminRoleName *string    `gorm:"-" json:"super_admin_role_name" validate:"omitempty,xss_clean" conform:"trim"`
	SuperAdminRoleId   *uuid.UUID `gorm:"-" json:"super_admin_role_id" validate:""`
	OrganizationRoleID int64      `json:"organization_role_id,omitempty" validate:""`
	Photo              string     `json:"photo" validate:"xss_clean" conform:"trim"`
	Address            string     `json:"address" validate:"xss_clean" conform:"trim"`
	LastLogin          *string    `gorm:"-" json:"last_login,omitempty" validate:""`
	Status             string     `json:"status" validate:"xss_clean" conform:"trim"`
	OpenstackData      string     `json:"openstack_data,omitempty" validate:"xss_clean" conform:"trim"`
	//
	VerifyToken         string            `json:"varify_token,omitempty" validate:"xss_clean" conform:"trim"`
	VerifedAt           string            `json:"verified_at" validate:"xss_clean" conform:"trim"`
	CreatedAt           time.Time         `gorm:"null" json:"created_at" validate:""`
	UpdatedAt           time.Time         `gorm:"null" json:"updated_at" validate:""`
	DeletedAt           gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty" validate:""`
	Organization        *Organization     `json:"organization" gorm:"foreignKey:OrganizationID" validate:""`
	OrganizationRole    *OrganizationRole `json:"organization_role" validate:""`
	IsCreator           bool              `json:"is_creator" validate:""`
	IsIdentityVerified  *bool             `gorm:"null, default: null" json:"is_identity_verified" validate:""`
	ReasonRejected      *string           `gorm:"null, default: null" json:"reason_rejected" validate:"omitempty,xss_clean" conform:"trim"`
	Project             []Project         `json:"project" gorm:"many2many:userprojects" validate:""`
	DekaprimeUserID     string            `gorm:"null" json:"dekaprime_user_id" validate:"xss_clean" conform:"trim"`
	IsMfaActive         bool              `gorm:"null" json:"is_mfa_active"`
	OnlyGoogleSSO       bool              `gorm:"null" json:"only_google_sso" gorm:"default:false"`
	ExpiredDatePassword time.Time         `gorm:"column:expired_date_password"`
}

// Role ..
type Role struct {
	ID        int64          `gorm:"primaryKey"`
	Name      string         `gorm:"not null" json:"name" `
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" json:"deleted_at"`
}

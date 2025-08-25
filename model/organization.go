package model

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID                         string         `gorm:"primaryKey;type:uuid" json:"id"`
	UserID                     string         `json:"user_id,omitempty"`
	SalesID                    *string        `gorm:"type:uuid;index" json:"sales_id"`
	Name                       string         `gorm:"not null" json:"name"`
	Prefix                     *string        `json:"prefix"`
	Segment                    *string        `json:"segment"`
	Subsegment                 *string        `json:"subsegment"`
	Email                      string         `gorm:"not null,unique" json:"email"`
	NpwpCorporate              string         `json:"npwp_corporate"`
	Details                    string         `json:"details"`
	PhoneNumber                string         `json:"phone_number_org"`
	Address                    string         `json:"address"`
	Pic                        *string        `json:"pic"`
	OpenstackPrefix            int64          `gorm:"auto_increment" json:"openstack_prefix"`
	Status                     int64          `gorm:"default:0" json:"status"`
	VClusterProjectName        string         `json:"v_cluster_project_name"`
	VClusterProjectNamespace   string         `json:"v_cluster_project_namespace"`
	CreatedAt                  time.Time      `gorm:"null" json:"created_at"`
	UpdatedAt                  time.Time      `gorm:"null" json:"updated_at"`
	DeletedAt                  gorm.DeletedAt `gorm:"index" json:",omitempty"`
	IsPersonal                 bool           `json:"is_personal"`
	SuspendType                *string        `gorm:"type:varchar(5)" json:"suspend_type"`
	DekaprimeOrganizationID    string         `gorm:"null" json:"dekaprime_organization_id"`
	ServiceActivationDate      *time.Time     `json:"service_activation_date"`
	SuspendDate                *time.Time     `json:"suspend_date"`
	TerminationDate            *time.Time     `json:"termination_date"`
	Country                    string         `json:"country"`
	City                       string         `json:"city"`
	Region                     string         `json:"region"`
	Zip                        string         `json:"zip"`
	PrefixBussinessPhoneNumber string         `json:"prefix_bussiness_phone_number"`
	PrefixPersonalPhoneNumber  string         `json:"prefix_personal_phone_number"`
	OrganizationCode           int64          `gorm:"auto_increment" json:"organization_code"`
	ExcludeReportTopups        bool           `json:"exclude_report_topups"`
	ClusterType                string         `gorm:"not null;default:'shared'" json:"cluster_type"`
}

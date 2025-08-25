package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID                        string         `gorm:"primaryKey" json:"id"`
	OrganizationID            string         `gorm:"not null;type:uuid" json:"organization_id,omitempty"`
	Name                      string         `gorm:"not null" json:"name"`
	OpenstackProjectID        string         `gorm:"not null" json:"openstack_project_id"`
	DekaprimeProjectID        string         `gorm:"null" json:"dekaprime_project_id"`
	Description               string         `json:"description"`
	Status                    int32          `gorm:"default:0" json:"status"`
	PaymentMethod             string         `gorm:"type:varchar(10);default:'prepaid'" json:"payment_method"`
	PostpaidType              *string        `gorm:"type:varchar(10)" json:"postpaid_type"`
	PostpaidFixedType         *string        `gorm:"null; type:varchar(10)" json:"postpaid_fixed_type"`
	ServiceActivationDate     *time.Time     `json:"service_activation_date"`
	NetworkNumber             *string        `json:"network_number"`
	IntervalBillingPerMonth   int            `gorm:"default: 1" json:"interval_billing_per_month"`
	InvoiceDuePeriodePerDay   int            `gorm:"default: 30" json:"invoice_due_periode_per_day"`
	FixedBillingPricePerMonth int            `gorm:"default: 0" json:"fixed_billing_price_per_month"`
	DocContract               *string        `json:"doc_contract"`
	Region                    string         `json:"region"`
	VclusterName              string         `json:"vcluster_name"`
	VclusterCluster           string         `json:"vcluster_cluster"`
	VatID                     string         `json:"vat_id"`
	VatName                   string         `json:"vat_name"`
	VatAddress                string         `json:"vat_address"`
	VclusterNamespace         string         `json:"vcluster_namespace"`
	ClusterRefNamespace       string         `json:"cluster_ref_namespace"`
	SuspendDate               *time.Time     `gorm:"null" json:"suspend_date"`
	TerminationDate           *time.Time     `gorm:"null" json:"termination_date"`
	CreatedAt                 time.Time      `gorm:"null" json:"created_at"`
	UpdatedAt                 time.Time      `gorm:"null" json:"updated_at"`
	DeletedAt                 gorm.DeletedAt `gorm:"index" json:",omitempty"`
	VclusterTemplate          string         `json:"vcluster_template"`
	IsPolicyProjectCreated    bool           `json:"is_policy_project_created"`
	IsSubnetProjectCreated    bool           `json:"is_subnet_project_created"`
	IsCiliumProjectCreated    bool           `json:"is_cilium_project_created"`
	IsCrd                     bool           `json:"is_crd"`
	IsDekaRespatiEnable       bool           `json:"is_deka_respati_enable" gorm:"default:false"`

	Organization *Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
}

func (Project) TableName() string {
	return "projects"
}

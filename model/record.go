package model

//type Record struct {
//	ID            string         `gorm:"primaryKey" json:"id"`
//	ZoneID        string         `gorm:"not null" json:"zoneId"`
//	Type          string         `gorm:"not null" json:"type"`
//	Status        bool           `gorm:"not null" json:"status"`
//	Cache         bool           `gorm:"not null" json:"cache"`
//	Name          string         `gorm:"not null" json:"name"`
//	Content       pq.StringArray `gorm:"not null;type:text[]" json:"content"`
//	TTL           uint32         `gorm:"ttl" json:"ttl"`
//	IsHidden      bool           `gorm:"not null" json:"isHidden"`
//	PricePerMonth float64        `gorm:"not null" json:"pricePerMonth"`
//	PricePerHour  float64        `gorm:"not null" json:"pricePerHour"`
//	CreatedBy     string         `gorm:"not null" json:"createdBy"`
//	UpdatedBy     string         `gorm:"not null" json:"updatedBy"`
//	CreatedAt     time.Time      `gorm:"not null" json:"createdAt"`
//	UpdatedAt     time.Time      `gorm:"not null" json:"updatedAt"`
//	DeletedAt     gorm.DeletedAt `gorm:"null" json:"-" `
//}
//
//func (m *Record) BeforeCreate(tx *gorm.DB) error {
//	now := time.Now()
//
//	if m.ID == "" {
//		m.ID = uuid.New().String()
//	}
//
//	if m.CreatedAt.IsZero() {
//		m.CreatedAt = now
//	}
//	if m.UpdatedAt.IsZero() {
//		m.UpdatedAt = now
//	}
//
//	return nil
//}
//
//func (m *Record) BeforeUpdate(tx *gorm.DB) (err error) {
//	now := time.Now()
//	m.UpdatedAt = now
//	return
//}
//
//func (Record) TableName() string {
//	return "dekadns.records"
//}

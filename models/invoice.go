package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CustomerID  uuid.UUID `gorm:"type:char(36);not null" json:"customer_id"`
	Customer    Customer  `gorm:"foreignKey:CustomerID" json:"customer"`
	UsageMonth  string    `gorm:"type:varchar(7);index" json:"usage_month"`
	UsageM3     float64   `json:"usage_m3"`
	Abonemen    float64   `json:"abonemen"`
	PricePerM3  float64   `json:"price_per_m3"`
	TotalAmount float64   `json:"total_amount"`
	IsPaid      bool      `gorm:"default:false" json:"is_paid"`
	TenantID    uuid.UUID `gorm:"type:char(36);index" json:"tenant_id"`
	CreatedAt   int64     `gorm:"autoCreateTime" json:"created_at"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.New()
	return
}

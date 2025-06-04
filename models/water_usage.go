package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WaterUsage struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CustomerID uuid.UUID `gorm:"type:char(36);not null;index" json:"customer_id"`
	Customer   Customer  `gorm:"foreignKey:CustomerID" json:"customer"`
	UsageMonth string    `gorm:"type:varchar(7);not null;index" json:"usage_month"` // e.g. 2025-06
	MeterStart float64   `json:"meter_start"`
	MeterEnd   float64   `json:"meter_end"`
	UsageM3    float64   `json:"usage_m3"`
	TenantID   uuid.UUID `gorm:"type:char(36);not null;index" json:"tenant_id"`
	CreatedAt  int64     `gorm:"autoCreateTime" json:"created_at"`
}

func (w *WaterUsage) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New()
	return
}

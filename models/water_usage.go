package models

import (
	"github.com/google/uuid"
)

type WaterUsage struct {
	CustomerID       uuid.UUID `gorm:"type:char(36);not null;index" json:"customer_id"`
	Customer         Customer  `gorm:"foreignKey:CustomerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"customer"`
	UsageMonth       string    `gorm:"type:varchar(7);not null;index" json:"usage_month"` // e.g. 2025-06
	MeterStart       float64   `json:"meter_start"`
	MeterEnd         float64   `json:"meter_end"`
	UsageM3          float64   `json:"usage_m3"`
	AmountCalculated float64   `json:"amount_calculated"` // hasil UsageM3 * tarif
	TenantID         uuid.UUID `gorm:"type:char(36);not null;index" json:"tenant_id"`

	BaseModel
}

// TableName overrides the table name for GORM
func (WaterUsage) TableName() string {
	return "water_usages"
}

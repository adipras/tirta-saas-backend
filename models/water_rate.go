package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WaterRate struct {
	ID             uuid.UUID        `gorm:"type:char(36);primaryKey" json:"id"`
	Amount         float64          `gorm:"not null" json:"amount"`
	EffectiveDate  time.Time        `gorm:"not null" json:"effective_date"`
	Active         bool             `gorm:"default:true" json:"active"`
	SubscriptionID uuid.UUID        `gorm:"type:char(36);not null" json:"subscription_id"`
	Subscription   SubscriptionType `gorm:"foreignKey:SubscriptionID" json:"subscription"`
	TenantID       uuid.UUID        `gorm:"type:char(36);not null;index" json:"tenant_id"`
}

func (w *WaterRate) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New()
	return
}

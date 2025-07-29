package models

import (
	"github.com/google/uuid"
)

type Customer struct {
	BaseModel

	MeterNumber    string           `gorm:"unique;not null" json:"meter_number"`
	Name           string           `gorm:"not null" json:"name"`
	Email          string           `gorm:"index" json:"email"`
	Password       string           `json:"-"`
	Address        string           `json:"address"`
	Phone          string           `json:"phone"`
	SubscriptionID uuid.UUID        `gorm:"type:char(36);not null" json:"subscription_id"`
	Subscription   SubscriptionType `gorm:"foreignKey:SubscriptionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"subscription"`
	IsActive       bool             `gorm:"default:false" json:"is_active"`
	TenantID       uuid.UUID        `gorm:"type:char(36);not null;index" json:"tenant_id"`
}

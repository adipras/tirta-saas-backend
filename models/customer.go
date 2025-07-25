package models

import (
	"github.com/google/uuid"
)

type Customer struct {
	BaseModel

	Name           string           `gorm:"not null" json:"name"`
	Address        string           `json:"address"`
	Phone          string           `json:"phone"`
	SubscriptionID uuid.UUID        `gorm:"type:char(36);not null" json:"subscription_id"`
	Subscription   SubscriptionType `gorm:"foreignKey:SubscriptionID" json:"subscription"`
	IsActive       bool             `gorm:"default:false" json:"is_active"`
	TenantID       uuid.UUID        `gorm:"type:char(36);not null;index" json:"tenant_id"`
}

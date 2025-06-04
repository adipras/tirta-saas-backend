package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID             uuid.UUID        `gorm:"type:char(36);primaryKey" json:"id"`
	Name           string           `gorm:"not null" json:"name"`
	Address        string           `json:"address"`
	Phone          string           `json:"phone"`
	SubscriptionID uuid.UUID        `gorm:"type:char(36);not null" json:"subscription_id"`
	Subscription   SubscriptionType `gorm:"foreignKey:SubscriptionID" json:"subscription"`
	TenantID       uuid.UUID        `gorm:"type:char(36);not null;index" json:"tenant_id"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}

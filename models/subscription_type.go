package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionType struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	TenantID    uuid.UUID `gorm:"type:char(36);not null;index" json:"tenant_id"`
}

func (s *SubscriptionType) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `json:"name"`
	VillageCode string    `gorm:"unique" json:"village_code"`
	CreatedAt   time.Time `json:"created_at"`
}

// Generate UUID otomatis saat insert
func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

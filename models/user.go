package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	TenantID  uuid.UUID `gorm:"type:char(36);not null" json:"tenant_id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`
	Role      string    `gorm:"type:enum('admin','operator');default:'operator'" json:"role"`
	CreatedAt time.Time `json:"created_at"`

	Tenant Tenant `gorm:"foreignKey:TenantID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

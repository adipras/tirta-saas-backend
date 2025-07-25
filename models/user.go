package models

import (
	"github.com/google/uuid"
)

type User struct {
	BaseModel

	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
	Role     string `gorm:"type:enum('admin','operator');default:'operator'" json:"role"`

	TenantID uuid.UUID `gorm:"type:char(36);not null" json:"tenant_id"`
	Tenant   Tenant    `gorm:"foreignKey:TenantID" json:"-"`
}

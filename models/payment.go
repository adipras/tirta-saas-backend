// models/payment.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	TenantID  uuid.UUID `gorm:"type:char(36);not null;index" json:"tenant_id"`
	InvoiceID uuid.UUID `gorm:"type:char(36);not null;index" json:"invoice_id"`
	Invoice   Invoice   `gorm:"foreignKey:InvoiceID" json:"invoice"`
	Amount    float64   `gorm:"not null" json:"amount"`
	Penalty   float64   `gorm:"default:0" json:"penalty"`
	PaidAt    time.Time `gorm:"not null" json:"paid_at"`

	BaseModel
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if err = p.BaseModel.BeforeCreate(tx); err != nil {
		return
	}
	p.PaidAt = time.Now()
	return
}

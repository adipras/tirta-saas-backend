// models/payment.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	InvoiceID uuid.UUID `gorm:"type:char(36);not null;index" json:"invoice_id"`
	Invoice   Invoice   `gorm:"foreignKey:InvoiceID"`
	Amount    float64   `gorm:"not null" json:"amount"`
	PaidAt    time.Time `gorm:"not null" json:"paid_at"`
	TenantID  uuid.UUID `gorm:"type:char(36);not null;index" json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	p.PaidAt = time.Now()
	return
}

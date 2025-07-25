package requests

import (
	"github.com/google/uuid"
)

type CreateWaterUsageRequest struct {
	CustomerID uuid.UUID `json:"customer_id" binding:"required"`
	UsageMonth string    `json:"usage_month" binding:"required,len=7"` // Format: YYYY-MM
	MeterEnd   float64   `json:"meter_end" binding:"required,gte=0"`
}

package requests

import "github.com/google/uuid"

type CreatePaymentRequest struct {
	InvoiceID uuid.UUID `json:"invoice_id" binding:"required"`
	Amount    float64   `json:"amount" binding:"required"`
}

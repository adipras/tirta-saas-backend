package responses

import "github.com/google/uuid"

type CustomerResponse struct {
	ID             uuid.UUID `json:"id"`
	MeterNumber    string    `json:"meter_number"`
	Name           string    `json:"name"`
	Email          string    `json:"email,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	Address        string    `json:"address,omitempty"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	IsActive       bool      `json:"is_active"`
}

type CustomerListResponse struct {
	Customers []CustomerResponse `json:"customers"`
	Total     int               `json:"total"`
}

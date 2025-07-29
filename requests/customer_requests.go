package requests

import "github.com/google/uuid"

type CreateCustomerRequest struct {
	MeterNumber    string    `json:"meter_number" binding:"required"`
	Name           string    `json:"name" binding:"required"`
	Email          string    `json:"email" binding:"required,email"`
	Password       string    `json:"password" binding:"required,min=6"`
	SubscriptionID uuid.UUID `json:"subscription_id" binding:"required"`
	Phone          string    `json:"phone,omitempty"`
	Address        string    `json:"address,omitempty"`
}

type UpdateCustomerRequest struct {
	Name           string    `json:"name" binding:"required"`
	SubscriptionID uuid.UUID `json:"subscription_id" binding:"required"`
	Phone          string    `json:"phone,omitempty"`
	Address        string    `json:"address,omitempty"`
}

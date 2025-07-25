package requests

import "github.com/google/uuid"

type CreateCustomerRequest struct {
	Name           string    `json:"name" binding:"required"`
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

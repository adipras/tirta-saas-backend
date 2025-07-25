package responses

import "github.com/google/uuid"

type CustomerResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Phone          string    `json:"phone,omitempty"`
	Address        string    `json:"address,omitempty"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	IsActive       bool      `json:"is_active"`
}

package requests

import "github.com/google/uuid"

type CreateTenantUserRequest struct {
	Name     string     `json:"name" binding:"required,min=3"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=6"`
	Role     string     `json:"role" binding:"required"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"` // Only for platform owners
}

type UpdateTenantUserRequest struct {
	Name string `json:"name,omitempty"`
	Role string `json:"role,omitempty"`
}
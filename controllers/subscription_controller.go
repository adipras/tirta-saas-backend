package controllers

import (
	"net/http"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"
	"github.com/adipras/tirta-saas-backend/requests"
	"github.com/adipras/tirta-saas-backend/responses"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateSubscriptionType(c *gin.Context) {
	var req requests.CreateSubscriptionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	sub := models.SubscriptionType{
		Name:            req.Name,
		Description:     req.Description,
		RegistrationFee: req.RegistrationFee,
		MonthlyFee:      req.MonthlyFee,
		MaintenanceFee:  req.MaintenanceFee,
		LateFeePerDay:   req.LateFeePerDay,
		MaxLateFee:      req.MaxLateFee,
		TenantID:        tenantID,
	}

	if err := config.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan SubscriptionType"})
		return
	}

	res := responses.SubscriptionTypeResponse{
		ID:              sub.ID,
		Name:            sub.Name,
		Description:     sub.Description,
		RegistrationFee: sub.RegistrationFee,
		MonthlyFee:      sub.MonthlyFee,
		MaintenanceFee:  sub.MaintenanceFee,
		LateFeePerDay:   sub.LateFeePerDay,
		MaxLateFee:      sub.MaxLateFee,
		CreatedAt:       sub.CreatedAt,
	}

	c.JSON(http.StatusCreated, res)
}

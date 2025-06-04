package controllers

import (
	"net/http"
	"time"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateWaterRate(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var input struct {
		Amount         float64   `json:"amount"`
		EffectiveDate  string    `json:"effective_date"` // YYYY-MM-DD
		SubscriptionID uuid.UUID `json:"subscription_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set existing rates for same subscription to inactive
	config.DB.Model(&models.WaterRate{}).
		Where("subscription_id = ? AND tenant_id = ?", input.SubscriptionID, tenantID).
		Update("active", false)

	date, err := time.Parse("2006-01-02", input.EffectiveDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal tidak valid"})
		return
	}

	rate := models.WaterRate{
		Amount:         input.Amount,
		EffectiveDate:  date,
		Active:         true,
		SubscriptionID: input.SubscriptionID,
		TenantID:       tenantID,
	}

	if err := config.DB.Create(&rate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat tarif"})
		return
	}

	c.JSON(http.StatusCreated, rate)
}

func GetWaterRates(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	var rates []models.WaterRate

	if err := config.DB.Preload("Subscription").
		Where("tenant_id = ?", tenantID).
		Order("effective_date DESC").
		Find(&rates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, rates)
}

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

func UpdateWaterRate(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	rateID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid water rate ID"})
		return
	}

	var rate models.WaterRate
	if err := config.DB.Where("id = ? AND tenant_id = ?", rateID, tenantID).
		First(&rate).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarif air tidak ditemukan"})
		return
	}

	var input struct {
		Amount        float64 `json:"amount"`
		EffectiveDate string  `json:"effective_date"` // YYYY-MM-DD
		Active        bool    `json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", input.EffectiveDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal tidak valid"})
		return
	}

	// If activating this rate, deactivate others for the same subscription
	if input.Active && !rate.Active {
		config.DB.Model(&models.WaterRate{}).
			Where("subscription_id = ? AND tenant_id = ? AND id != ?", 
				rate.SubscriptionID, tenantID, rateID).
			Update("active", false)
	}

	rate.Amount = input.Amount
	rate.EffectiveDate = date
	rate.Active = input.Active

	if err := config.DB.Save(&rate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui tarif"})
		return
	}

	c.JSON(http.StatusOK, rate)
}

func DeleteWaterRate(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	rateID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid water rate ID"})
		return
	}

	var rate models.WaterRate
	if err := config.DB.Where("id = ? AND tenant_id = ?", rateID, tenantID).
		First(&rate).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarif air tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&rate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus tarif"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tarif air berhasil dihapus"})
}

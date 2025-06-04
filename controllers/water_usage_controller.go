package controllers

import (
	"net/http"
	"time"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateWaterUsage(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var input struct {
		CustomerID uuid.UUID `json:"customer_id"`
		UsageMonth string    `json:"usage_month"` // format: YYYY-MM
		MeterEnd   float64   `json:"meter_end"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hitung bulan sebelumnya
	prevMonth, err := time.Parse("2006-01", input.UsageMonth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format bulan tidak valid. Gunakan YYYY-MM"})
		return
	}
	prevMonth = prevMonth.AddDate(0, -1, 0)
	prevMonthStr := prevMonth.Format("2006-01")

	// Ambil meter_end bulan sebelumnya
	var lastUsage models.WaterUsage
	meterStart := 0.0
	if err := config.DB.Where("customer_id = ? AND usage_month = ? AND tenant_id = ?", input.CustomerID, prevMonthStr, tenantID).
		First(&lastUsage).Error; err == nil {
		meterStart = lastUsage.MeterEnd
	}

	if input.MeterEnd < meterStart {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meter akhir lebih kecil dari meter sebelumnya"})
		return
	}

	usage := models.WaterUsage{
		CustomerID: input.CustomerID,
		UsageMonth: input.UsageMonth,
		MeterStart: meterStart,
		MeterEnd:   input.MeterEnd,
		UsageM3:    input.MeterEnd - meterStart,
		TenantID:   tenantID,
	}

	if err := config.DB.Create(&usage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusCreated, usage)
}

func GetWaterUsages(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	var records []models.WaterUsage

	if err := config.DB.Preload("Customer").
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, records)
}

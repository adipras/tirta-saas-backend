package controllers

import (
	"net/http"
	"time"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"
	"github.com/adipras/tirta-saas-backend/requests"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateWaterUsage(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var req requests.CreateWaterUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hitung bulan sebelumnya
	prevMonth, err := time.Parse("2006-01", req.UsageMonth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format bulan tidak valid. Gunakan YYYY-MM"})
		return
	}
	prevMonth = prevMonth.AddDate(0, -1, 0)
	prevMonthStr := prevMonth.Format("2006-01")

	// Ambil meter_end bulan sebelumnya
	var lastUsage models.WaterUsage
	meterStart := 0.0
	if err := config.DB.Where("customer_id = ? AND usage_month = ? AND tenant_id = ?", req.CustomerID, prevMonthStr, tenantID).
		First(&lastUsage).Error; err == nil {
		meterStart = lastUsage.MeterEnd
	}

	if req.MeterEnd < meterStart {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meter akhir lebih kecil dari meter sebelumnya"})
		return
	}

	// Ambil data customer
	var customer models.Customer
	if err := config.DB.Where("id = ?", req.CustomerID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	// Ambil tarif aktif untuk subscription pelanggan
	var rate models.WaterRate
	if err := config.DB.
		Where("subscription_id = ? AND active = ?", customer.SubscriptionID, true).
		Order("effective_date DESC").
		First(&rate).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tarif air aktif tidak ditemukan"})
		return
	}

	UsageM3 := req.MeterEnd - meterStart
	usage := models.WaterUsage{
		CustomerID:       req.CustomerID,
		UsageMonth:       req.UsageMonth,
		MeterStart:       meterStart,
		MeterEnd:         req.MeterEnd,
		UsageM3:          UsageM3,
		AmountCalculated: UsageM3 * rate.Amount,
		TenantID:         tenantID,
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

func GetWaterUsageByID(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var usage models.WaterUsage
	if err := config.DB.Preload("Customer").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&usage).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, usage)
}

func UpdateWaterUsage(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var input struct {
		MeterEnd float64 `json:"meter_end"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var usage models.WaterUsage
	if err := config.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&usage).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	if input.MeterEnd < usage.MeterStart {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meter akhir tidak boleh lebih kecil dari awal"})
		return
	}

	// Ambil data customer
	var customer models.Customer
	if err := config.DB.Where("id = ?", usage.CustomerID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	// Ambil tarif aktif untuk subscription pelanggan
	var rate models.WaterRate
	if err := config.DB.
		Where("subscription_id = ? AND active = ?", customer.SubscriptionID, true).
		Order("effective_date DESC").
		First(&rate).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tarif air aktif tidak ditemukan"})
		return
	}

	UsageM3 := input.MeterEnd - usage.MeterStart

	usage.MeterEnd = input.MeterEnd
	usage.UsageM3 = UsageM3
	usage.AmountCalculated = UsageM3 * rate.Amount

	if err := config.DB.Save(&usage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, usage)
}

func DeleteWaterUsage(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var usage models.WaterUsage
	if err := config.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&usage).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&usage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

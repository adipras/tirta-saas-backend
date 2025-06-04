package controllers

import (
	"net/http"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateInvoice(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var input struct {
		CustomerID uuid.UUID `json:"customer_id"`
		UsageMonth string    `json:"usage_month"` // e.g. 2025-06
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah invoice sudah dibuat
	var existing models.Invoice
	if err := config.DB.Where("customer_id = ? AND usage_month = ? AND tenant_id = ?", input.CustomerID, input.UsageMonth, tenantID).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tagihan untuk bulan ini sudah ada"})
		return
	}

	// Ambil data pelanggan terlebih dahulu (termasuk Subscription)
	var customer models.Customer
	if err := config.DB.Preload("Subscription").
		Where("id = ? AND tenant_id = ?", input.CustomerID, tenantID).
		First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	// Ambil data pemakaian air
	var usage models.WaterUsage
	if err := config.DB.Where("customer_id = ? AND usage_month = ? AND tenant_id = ?", input.CustomerID, input.UsageMonth, tenantID).
		First(&usage).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pemakaian air tidak ditemukan"})
		return
	}

	// Ambil tarif aktif berdasarkan subscription dan tenant
	var rate models.WaterRate
	if err := config.DB.
		Where("subscription_id = ? AND tenant_id = ? AND active = ?", customer.SubscriptionID, tenantID, true).
		Order("effective_date desc").
		First(&rate).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarif air aktif tidak ditemukan untuk pelanggan ini"})
		return
	}

	// Ambil abonemen dari subscription
	abonemen := customer.Subscription.SubscriptionFee

	invoice := models.Invoice{
		CustomerID:  input.CustomerID,
		UsageMonth:  input.UsageMonth,
		UsageM3:     usage.UsageM3,
		Abonemen:    abonemen,
		PricePerM3:  rate.Amount,
		TotalAmount: abonemen + (rate.Amount * usage.UsageM3),
		TenantID:    tenantID,
	}

	if err := config.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat tagihan"})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

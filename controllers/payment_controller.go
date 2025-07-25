package controllers

import (
	"fmt"
	"net/http"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"
	"github.com/adipras/tirta-saas-backend/requests"
	"github.com/adipras/tirta-saas-backend/responses"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePayment(c *gin.Context) {
	var req requests.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	// Ambil invoice terkait
	var invoice models.Invoice
	if err := config.DB.Where("id = ? AND tenant_id = ?", req.InvoiceID, tenantID).
		First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice tidak ditemukan"})
		return
	}

	if invoice.IsPaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tagihan sudah lunas"})
		return
	}

	if invoice.TotalPaid+req.Amount > invoice.TotalAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Pembayaran melebihi total tagihan. Sisa tagihan: %.2f", invoice.TotalAmount-invoice.TotalPaid),
		})
		return
	}

	// Buat record pembayaran
	payment := models.Payment{
		InvoiceID: req.InvoiceID,
		Amount:    req.Amount,
		TenantID:  tenantID,
	}
	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat pembayaran"})
		return
	}

	// Hitung total bayar baru
	var totalPaid float64
	config.DB.Model(&models.Payment{}).
		Where("invoice_id = ?", req.InvoiceID).
		Select("SUM(amount)").Scan(&totalPaid)

	// Update invoice
	invoice.TotalPaid = totalPaid
	invoice.IsPaid = totalPaid >= invoice.TotalAmount
	if err := config.DB.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status invoice"})
		return
	}

	// Jika invoice pendaftaran dan sudah lunas → aktifkan customer
	if invoice.Type == "registration" && invoice.IsPaid {
		if err := config.DB.Model(&models.Customer{}).
			Where("id = ? AND tenant_id = ?", invoice.CustomerID, tenantID).
			Update("is_active", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengaktifkan pelanggan"})
			return
		}
	}

	// Kirim response
	res := responses.PaymentResponse{
		ID:        payment.ID,
		InvoiceID: payment.InvoiceID,
		Amount:    payment.Amount,
		PaidAt:    payment.CreatedAt,
	}
	c.JSON(http.StatusCreated, res)
}

func GetPaymentHistoryByCustomerID(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	customerIDStr := c.Param("customer_id")

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id tidak valid"})
		return
	}

	var payments []models.Payment
	if err := config.DB.Preload("Invoice").
		Where("tenant_id = ? AND invoice_id IN (SELECT id FROM invoices WHERE customer_id = ?)", tenantID, customerID).
		Order("paid_at desc").
		Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat pembayaran"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

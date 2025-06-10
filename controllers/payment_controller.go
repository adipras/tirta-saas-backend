// controllers/payment_controller.go
package controllers

import (
	"fmt"
	"net/http"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePayment(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var input struct {
		InvoiceID uuid.UUID `json:"invoice_id"`
		Amount    float64   `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil invoice terkait
	var invoice models.Invoice
	if err := config.DB.Where("id = ? AND tenant_id = ?", input.InvoiceID, tenantID).
		First(&invoice).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice tidak ditemukan"})
		return
	}

	if invoice.IsPaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tagihan sudah lunas"})
		return
	}

	if invoice.TotalPaid+input.Amount > invoice.TotalAmount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Pembayaran melebihi total tagihan. Sisa tagihan: %.2f", invoice.TotalAmount-invoice.TotalPaid),
		})
		return
	}

	// Buat record pembayaran
	payment := models.Payment{
		InvoiceID: input.InvoiceID,
		Amount:    input.Amount,
		TenantID:  tenantID,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat pembayaran"})
		return
	}

	// Hitung total bayar
	var totalPaid float64
	config.DB.Model(&models.Payment{}).
		Where("invoice_id = ?", input.InvoiceID).
		Select("SUM(amount)").Scan(&totalPaid)

	// Update invoice
	invoice.TotalPaid = totalPaid
	invoice.IsPaid = totalPaid >= invoice.TotalAmount

	if err := config.DB.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status invoice"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"payment": payment,
		"invoice": invoice,
	})
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

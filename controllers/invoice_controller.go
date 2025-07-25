package controllers

import (
	"net/http"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateMonthlyInvoice(c *gin.Context) {
	type Request struct {
		UsageMonth string `json:"usage_month" binding:"required"` // format: YYYY-MM
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UsageMonth wajib diisi (format: YYYY-MM)"})
		return
	}

	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	// Ambil semua WaterUsage bulan tsb yang belum dibuatkan Invoice
	var usages []models.WaterUsage
	if err := config.DB.
		Where("usage_month = ? AND tenant_id = ?", req.UsageMonth, tenantID).
		Find(&usages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal ambil data water usage"})
		return
	}

	if len(usages) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Tidak ada water usage untuk bulan tersebut"})
		return
	}

	created := 0
	skipped := 0

	for _, usage := range usages {
		// Cek apakah invoice sudah pernah dibuat
		var existing models.Invoice
		err := config.DB.Where("customer_id = ? AND usage_month = ? AND type = ?",
			usage.CustomerID, usage.UsageMonth, "monthly").First(&existing).Error
		if err == nil {
			skipped++
			continue
		}

		// Ambil data pelanggan & SubscriptionType
		var customer models.Customer
		if err := config.DB.First(&customer, "id = ?", usage.CustomerID).Error; err != nil {
			continue
		}

		var subType models.SubscriptionType
		if err := config.DB.First(&subType, "id = ?", customer.SubscriptionID).Error; err != nil {
			continue
		}

		// // Hitung denda dari tagihan bulan sebelumnya jika ada dan belum dibayar
		// var penalty float64 = 0
		// prevMonth, _ := utils.PreviousMonth(usage.UsageMonth)
		// var prevInvoice models.Invoice
		// err = config.DB.Where("customer_id = ? AND usage_month = ? AND type = ?", usage.CustomerID, prevMonth, "monthly").
		// 	First(&prevInvoice).Error

		// if err == nil && !prevInvoice.IsPaid {
		// 	// Hitung keterlambatan
		// 	// Misal batas pembayaran adalah 10 hari setelah akhir bulan penggunaan
		// 	dueDate, _ := utils.DueDateFromUsageMonth(prevInvoice.UsageMonth, 10) // 2025-05 → 2025-06-10
		// 	today := time.Now()
		// 	if today.After(dueDate) {
		// 		daysLate := int(today.Sub(dueDate).Hours() / 24)
		// 		penalty = float64(daysLate) * subType.LateFeePerDay
		// 		if penalty > subType.MaxLateFee {
		// 			penalty = subType.MaxLateFee
		// 		}
		// 	}
		// }

		total := usage.AmountCalculated + subType.MonthlyFee + subType.MaintenanceFee

		invoice := models.Invoice{
			CustomerID:  usage.CustomerID,
			UsageMonth:  usage.UsageMonth,
			UsageM3:     usage.UsageM3,
			Abonemen:    subType.MonthlyFee,
			PricePerM3:  usage.AmountCalculated / usage.UsageM3,
			TotalAmount: total,
			TotalPaid:   0,
			IsPaid:      false,
			TenantID:    tenantID,
			Type:        "monthly",
		}

		if err := config.DB.Create(&invoice).Error; err == nil {
			created++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Generate invoice selesai",
		"created_count": created,
		"skipped":       skipped,
	})
}

func GetInvoices(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	var invoices []models.Invoice

	if err := config.DB.Preload("Customer").
		Where("tenant_id = ?", tenantID).
		Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

// func UpdateInvoice(c *gin.Context) {
// 	tenantID := c.MustGet("tenant_id").(uuid.UUID)
// 	id := c.Param("id")

// 	var invoice models.Invoice
// 	if err := config.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
// 		First(&invoice).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice tidak ditemukan"})
// 		return
// 	}

// 	var input dto.UpdateInvoiceInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	invoice.UsageM3 = input.UsageM3
// 	invoice.Abonemen = input.Abonemen
// 	invoice.PricePerM3 = input.PricePerM3
// 	invoice.TotalAmount = input.TotalAmount
// 	invoice.IsPaid = input.IsPaid
// 	invoice.TotalPaid = input.TotalPaid

// 	if err := config.DB.Save(&invoice).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui invoice"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, invoice)
// }

func DeleteInvoice(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	if err := config.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.Invoice{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice berhasil dihapus"})
}

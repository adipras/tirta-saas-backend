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

func CreateCustomer(c *gin.Context) {
	var req requests.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantIDValue, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	tenantID := tenantIDValue.(uuid.UUID)

	// Ambil SubscriptionType
	var subType models.SubscriptionType
	if err := config.DB.Where("id = ? AND tenant_id = ?", req.SubscriptionID, tenantID).First(&subType).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription type not found"})
		return
	}

	// Buat Customer
	customer := models.Customer{
		Name:           req.Name,
		SubscriptionID: req.SubscriptionID,
		IsActive:       false,
		TenantID:       tenantID,
	}
	if err := config.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	// Buat Invoice untuk biaya pendaftaran
	invoice := models.Invoice{
		CustomerID:  customer.ID,
		UsageMonth:  "", // Kosong karena ini bukan invoice pemakaian
		UsageM3:     0,
		Abonemen:    0,
		PricePerM3:  0,
		TotalAmount: subType.RegistrationFee,
		IsPaid:      false,
		TotalPaid:   0,
		Type:        "registration",
		TenantID:    tenantID,
	}
	if err := config.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create registration invoice"})
		return
	}

	// Respon
	response := responses.CustomerResponse{
		ID:             customer.ID,
		Name:           customer.Name,
		Address:        customer.Address,
		Phone:          customer.Phone,
		SubscriptionID: customer.SubscriptionID,
		IsActive:       customer.IsActive,
	}
	c.JSON(http.StatusCreated, response)
}

func GetCustomers(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	var customers []models.Customer

	if err := config.DB.Preload("Subscription").
		Where("tenant_id = ?", tenantID).
		Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func UpdateCustomer(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	var customer models.Customer
	if err := config.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	var input requests.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.SubscriptionID = input.SubscriptionID

	if err := config.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui pelanggan"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func DeleteCustomer(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	if err := config.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Customer{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pelanggan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pelanggan berhasil dihapus"})
}

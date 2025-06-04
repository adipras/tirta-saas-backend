package controllers

import (
	"net/http"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCustomer(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.TenantID = tenantID

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pelanggan"})
		return
	}

	c.JSON(http.StatusCreated, input)
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

	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer.Name = input.Name
	customer.Address = input.Address
	customer.Phone = input.Phone
	customer.SubscriptionID = input.SubscriptionID

	config.DB.Save(&customer)

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

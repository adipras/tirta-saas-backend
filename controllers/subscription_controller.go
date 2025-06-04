package controllers

import (
	"net/http"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateSubscriptionType(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var input models.SubscriptionType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.TenantID = tenantID

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat abonemen"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func GetAllSubscriptionTypes(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	var data []models.SubscriptionType

	if err := config.DB.Where("tenant_id = ?", tenantID).Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func UpdateSubscriptionType(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	var data models.SubscriptionType
	if err := config.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&data).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	var input models.SubscriptionType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data.Name = input.Name
	data.Description = input.Description

	config.DB.Save(&data)

	c.JSON(http.StatusOK, data)
}

func DeleteSubscriptionType(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	if err := config.DB.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.SubscriptionType{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil dihapus"})
}

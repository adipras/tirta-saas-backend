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

func CreateSubscriptionType(c *gin.Context) {
	var req requests.CreateSubscriptionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	sub := models.SubscriptionType{
		Name:            req.Name,
		Description:     req.Description,
		RegistrationFee: req.RegistrationFee,
		MonthlyFee:      req.MonthlyFee,
		MaintenanceFee:  req.MaintenanceFee,
		LateFeePerDay:   req.LateFeePerDay,
		MaxLateFee:      req.MaxLateFee,
		TenantID:        tenantID,
	}

	if err := config.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan SubscriptionType"})
		return
	}

	res := responses.SubscriptionTypeResponse{
		ID:              sub.ID,
		Name:            sub.Name,
		Description:     sub.Description,
		RegistrationFee: sub.RegistrationFee,
		MonthlyFee:      sub.MonthlyFee,
		MaintenanceFee:  sub.MaintenanceFee,
		LateFeePerDay:   sub.LateFeePerDay,
		MaxLateFee:      sub.MaxLateFee,
		CreatedAt:       sub.CreatedAt,
	}

	c.JSON(http.StatusCreated, res)
}

func GetAllSubscriptionTypes(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)

	var subscriptions []models.SubscriptionType
	if err := config.DB.Where("tenant_id = ?", tenantID).Find(&subscriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data subscription types"})
		return
	}

	var responseList []responses.SubscriptionTypeResponse
	for _, sub := range subscriptions {
		res := responses.SubscriptionTypeResponse{
			ID:              sub.ID,
			Name:            sub.Name,
			Description:     sub.Description,
			RegistrationFee: sub.RegistrationFee,
			MonthlyFee:      sub.MonthlyFee,
			MaintenanceFee:  sub.MaintenanceFee,
			LateFeePerDay:   sub.LateFeePerDay,
			MaxLateFee:      sub.MaxLateFee,
			CreatedAt:       sub.CreatedAt,
		}
		responseList = append(responseList, res)
	}

	c.JSON(http.StatusOK, responseList)
}

func GetSubscriptionType(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	subID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription type ID"})
		return
	}

	var sub models.SubscriptionType
	if err := config.DB.Where("id = ? AND tenant_id = ?", subID, tenantID).First(&sub).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription type tidak ditemukan"})
		return
	}

	res := responses.SubscriptionTypeResponse{
		ID:              sub.ID,
		Name:            sub.Name,
		Description:     sub.Description,
		RegistrationFee: sub.RegistrationFee,
		MonthlyFee:      sub.MonthlyFee,
		MaintenanceFee:  sub.MaintenanceFee,
		LateFeePerDay:   sub.LateFeePerDay,
		MaxLateFee:      sub.MaxLateFee,
		CreatedAt:       sub.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func UpdateSubscriptionType(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	subID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription type ID"})
		return
	}

	var req requests.CreateSubscriptionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sub models.SubscriptionType
	if err := config.DB.Where("id = ? AND tenant_id = ?", subID, tenantID).First(&sub).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription type tidak ditemukan"})
		return
	}

	sub.Name = req.Name
	sub.Description = req.Description
	sub.RegistrationFee = req.RegistrationFee
	sub.MonthlyFee = req.MonthlyFee
	sub.MaintenanceFee = req.MaintenanceFee
	sub.LateFeePerDay = req.LateFeePerDay
	sub.MaxLateFee = req.MaxLateFee

	if err := config.DB.Save(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui subscription type"})
		return
	}

	res := responses.SubscriptionTypeResponse{
		ID:              sub.ID,
		Name:            sub.Name,
		Description:     sub.Description,
		RegistrationFee: sub.RegistrationFee,
		MonthlyFee:      sub.MonthlyFee,
		MaintenanceFee:  sub.MaintenanceFee,
		LateFeePerDay:   sub.LateFeePerDay,
		MaxLateFee:      sub.MaxLateFee,
		CreatedAt:       sub.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

func DeleteSubscriptionType(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	id := c.Param("id")

	subID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription type ID"})
		return
	}

	var sub models.SubscriptionType
	if err := config.DB.Where("id = ? AND tenant_id = ?", subID, tenantID).First(&sub).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription type tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus subscription type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription type berhasil dihapus"})
}

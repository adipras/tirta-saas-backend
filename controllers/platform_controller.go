package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/models"
	"github.com/adipras/tirta-saas-backend/requests"
	"github.com/adipras/tirta-saas-backend/responses"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListTenants lists all tenants with pagination and filters (Platform Owner only)
func ListTenants(c *gin.Context) {
	var req requests.TenantSearchRequest
	
	// Set defaults
	if c.Query("page") == "" {
		req.Page = 1
	}
	if c.Query("page_size") == "" {
		req.PageSize = 20
	}
	if c.Query("sort_order") == "" {
		req.SortOrder = "desc"
	}
	if c.Query("sort_by") == "" {
		req.SortBy = "created_at"
	}
	
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}
	
	var tenants []models.Tenant
	query := config.DB.Model(&models.Tenant{})
	
	// Apply filters
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR village_code LIKE ? OR email LIKE ?", 
			searchPattern, searchPattern, searchPattern)
	}
	
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	
	if req.SubscriptionPlan != "" {
		query = query.Where("subscription_plan = ?", req.SubscriptionPlan)
	}
	
	// Count total records
	var total int64
	query.Count(&total)
	
	// Apply sorting
	sortField := req.SortBy
	if sortField == "" {
		sortField = "created_at"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortField, req.SortOrder))
	
	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)
	
	if err := query.Find(&tenants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch tenants",
			Error:   err.Error(),
		})
		return
	}
	
	// Transform to response
	var tenantList []responses.TenantListResponse
	for _, tenant := range tenants {
		tenantList = append(tenantList, responses.TenantListResponse{
			ID:                 tenant.ID,
			Name:               tenant.Name,
			VillageCode:        tenant.VillageCode,
			Email:              tenant.Email,
			Phone:              tenant.Phone,
			Status:             string(tenant.Status),
			SubscriptionPlan:   tenant.SubscriptionPlan,
			SubscriptionStatus: tenant.SubscriptionStatus,
			SubscriptionEndsAt: tenant.SubscriptionEndsAt,
			TotalUsers:         tenant.TotalUsers,
			TotalCustomers:     tenant.TotalCustomers,
			StorageUsedGB:      tenant.StorageUsedGB,
			CreatedAt:          tenant.CreatedAt,
		})
	}
	
	c.JSON(http.StatusOK, responses.PaginatedResponse{
		Status:  "success",
		Message: "Tenants retrieved successfully",
		Data:    tenantList,
		Meta: responses.PaginationMeta{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
			TotalPages:  int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
			TotalItems:  int(total),
		},
	})
}

// GetTenantDetail gets detailed tenant information (Platform Owner only)
func GetTenantDetail(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid tenant ID format",
		})
		return
	}
	
	var tenant models.Tenant
	if err := config.DB.First(&tenant, "id = ?", tenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  "error",
			Message: "Tenant not found",
		})
		return
	}
	
	response := responses.TenantDetailResponse{
		ID:                 tenant.ID,
		Name:               tenant.Name,
		VillageCode:        tenant.VillageCode,
		Email:              tenant.Email,
		Phone:              tenant.Phone,
		Address:            tenant.Address,
		Status:             string(tenant.Status),
		SubscriptionPlan:   tenant.SubscriptionPlan,
		SubscriptionStatus: tenant.SubscriptionStatus,
		SubscriptionEndsAt: tenant.SubscriptionEndsAt,
		TotalUsers:         tenant.TotalUsers,
		TotalCustomers:     tenant.TotalCustomers,
		StorageUsedGB:      tenant.StorageUsedGB,
		SuspendedAt:        tenant.SuspendedAt,
		SuspensionReason:   tenant.SuspensionReason,
		Notes:              tenant.Notes,
		CreatedAt:          tenant.CreatedAt,
		UpdatedAt:          tenant.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Tenant details retrieved successfully",
		Data:    response,
	})
}

// UpdateTenant updates tenant information (Platform Owner only)
func UpdateTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid tenant ID format",
		})
		return
	}
	
	var req requests.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	
	var tenant models.Tenant
	if err := config.DB.First(&tenant, "id = ?", tenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  "error",
			Message: "Tenant not found",
		})
		return
	}
	
	// Update fields
	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Email != "" {
		tenant.Email = req.Email
	}
	if req.Phone != "" {
		tenant.Phone = req.Phone
	}
	if req.Address != "" {
		tenant.Address = req.Address
	}
	if req.Notes != "" {
		tenant.Notes = req.Notes
	}
	if req.SubscriptionPlan != "" {
		tenant.SubscriptionPlan = req.SubscriptionPlan
	}
	
	if err := config.DB.Save(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  "error",
			Message: "Failed to update tenant",
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Tenant updated successfully",
		Data:    tenant,
	})
}

// SuspendTenant suspends a tenant (Platform Owner only)
func SuspendTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid tenant ID format",
		})
		return
	}
	
	var req requests.SuspendTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	
	var tenant models.Tenant
	if err := config.DB.First(&tenant, "id = ?", tenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  "error",
			Message: "Tenant not found",
		})
		return
	}
	
	if tenant.Status == models.TenantStatusSuspended {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Tenant is already suspended",
		})
		return
	}
	
	now := time.Now()
	tenant.Status = models.TenantStatusSuspended
	tenant.SuspendedAt = &now
	tenant.SuspensionReason = req.Reason
	
	if err := config.DB.Save(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  "error",
			Message: "Failed to suspend tenant",
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Tenant suspended successfully",
		Data:    tenant,
	})
}

// ActivateTenant activates a suspended tenant (Platform Owner only)
func ActivateTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid tenant ID format",
		})
		return
	}
	
	var tenant models.Tenant
	if err := config.DB.First(&tenant, "id = ?", tenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  "error",
			Message: "Tenant not found",
		})
		return
	}
	
	if tenant.Status == models.TenantStatusActive {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Tenant is already active",
		})
		return
	}
	
	tenant.Status = models.TenantStatusActive
	tenant.SuspendedAt = nil
	tenant.SuspensionReason = ""
	
	if err := config.DB.Save(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  "error",
			Message: "Failed to activate tenant",
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Tenant activated successfully",
		Data:    tenant,
	})
}

// DeleteTenant soft deletes a tenant (Platform Owner only)
func DeleteTenant(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid tenant ID format",
		})
		return
	}
	
	var tenant models.Tenant
	if err := config.DB.First(&tenant, "id = ?", tenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  "error",
			Message: "Tenant not found",
		})
		return
	}
	
	// Soft delete
	if err := config.DB.Delete(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete tenant",
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Tenant deleted successfully",
	})
}

// GetTenantStatistics gets tenant usage statistics (Platform Owner only)
func GetTenantStatistics(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid tenant ID format",
		})
		return
	}
	
	var tenant models.Tenant
	if err := config.DB.First(&tenant, "id = ?", tenantID).Error; err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Status:  "error",
			Message: "Tenant not found",
		})
		return
	}
	
	stats := responses.TenantStatisticsResponse{
		TenantID:   tenant.ID,
		TenantName: tenant.Name,
	}
	
	// Count users
	var totalUsers, activeUsers int64
	config.DB.Model(&models.User{}).Where("tenant_id = ?", tenantID).Count(&totalUsers)
	config.DB.Model(&models.User{}).Where("tenant_id = ? AND is_active = ?", tenantID, true).Count(&activeUsers)
	stats.TotalUsers = int(totalUsers)
	stats.ActiveUsers = int(activeUsers)
	
	// Count customers
	var totalCustomers, activeCustomers int64
	config.DB.Model(&models.Customer{}).Where("tenant_id = ?", tenantID).Count(&totalCustomers)
	config.DB.Model(&models.Customer{}).Where("tenant_id = ? AND is_active = ?", tenantID, true).Count(&activeCustomers)
	stats.TotalCustomers = int(totalCustomers)
	stats.ActiveCustomers = int(activeCustomers)
	stats.InactiveCustomers = stats.TotalCustomers - stats.ActiveCustomers
	
	// Invoice statistics
	var totalInvoices, paidInvoices, unpaidInvoices int64
	config.DB.Model(&models.Invoice{}).Where("tenant_id = ?", tenantID).Count(&totalInvoices)
	config.DB.Model(&models.Invoice{}).Where("tenant_id = ? AND payment_status = ?", tenantID, "PAID").Count(&paidInvoices)
	config.DB.Model(&models.Invoice{}).Where("tenant_id = ? AND payment_status != ?", tenantID, "PAID").Count(&unpaidInvoices)
	stats.TotalInvoices = int(totalInvoices)
	stats.PaidInvoices = int(paidInvoices)
	stats.UnpaidInvoices = int(unpaidInvoices)
	
	// Revenue statistics
	var totalRevenue, outstandingAmount float64
	config.DB.Model(&models.Payment{}).Where("tenant_id = ?", tenantID).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)
	config.DB.Model(&models.Invoice{}).Where("tenant_id = ? AND payment_status != ?", tenantID, "PAID").
		Select("COALESCE(SUM(total_amount - paid_amount), 0)").Scan(&outstandingAmount)
	
	stats.TotalRevenue = totalRevenue
	stats.OutstandingAmount = outstandingAmount
	
	// Water usage statistics
	var totalUsage float64
	config.DB.Model(&models.WaterUsage{}).Where("tenant_id = ?", tenantID).
		Select("COALESCE(SUM(usage_m3), 0)").Scan(&totalUsage)
	stats.TotalWaterUsage = totalUsage
	
	if stats.TotalCustomers > 0 {
		stats.AverageUsagePerCustomer = totalUsage / float64(stats.TotalCustomers)
	}
	
	// Storage and limits from subscription
	var subscription models.TenantSubscription
	if err := config.DB.Where("tenant_id = ? AND status = ?", tenantID, "ACTIVE").First(&subscription).Error; err == nil {
		stats.StorageLimitGB = subscription.MaxStorageGB
		stats.APICallsLimit = subscription.MaxAPICallsPerDay
	}
	
	stats.StorageUsedGB = tenant.StorageUsedGB
	stats.APICallsToday = 0 // TODO: Implement from metrics
	
	// Last activity
	var lastLog models.AuditLog
	if err := config.DB.Where("tenant_id = ?", tenantID).Order("created_at DESC").First(&lastLog).Error; err == nil {
		stats.LastActivityAt = &lastLog.CreatedAt
	}
	
	stats.StatisticsAsOf = time.Now()
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Tenant statistics retrieved successfully",
		Data:    stats,
	})
}

// GetPlatformAnalyticsOverview gets platform-wide overview statistics (Platform Owner only)
func GetPlatformAnalyticsOverview(c *gin.Context) {
	var stats responses.PlatformAnalyticsOverviewResponse
	
	// Tenant statistics
	var totalTenants, activeTenants, suspendedTenants, trialTenants int64
	config.DB.Model(&models.Tenant{}).Count(&totalTenants)
	config.DB.Model(&models.Tenant{}).Where("status = ?", "ACTIVE").Count(&activeTenants)
	config.DB.Model(&models.Tenant{}).Where("status = ?", "SUSPENDED").Count(&suspendedTenants)
	stats.TotalTenants = int(totalTenants)
	stats.ActiveTenants = int(activeTenants)
	stats.SuspendedTenants = int(suspendedTenants)
	
	// Trial tenants (subscriptions with trial)
	config.DB.Model(&models.TenantSubscription{}).
		Where("status = ? AND trial_ends_at > ?", "TRIAL", time.Now()).Count(&trialTenants)
	stats.TrialTenants = int(trialTenants)
	
	// Revenue statistics - total all time
	var totalRevenue float64
	config.DB.Model(&models.Payment{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)
	stats.TotalRevenue = totalRevenue
	
	// Monthly revenue - current month
	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	var monthlyRevenue float64
	config.DB.Model(&models.Payment{}).Where("created_at >= ?", firstDayOfMonth).
		Select("COALESCE(SUM(amount), 0)").Scan(&monthlyRevenue)
	stats.MonthlyRevenue = monthlyRevenue
	
	// Outstanding revenue
	var outstandingRevenue float64
	config.DB.Model(&models.Invoice{}).Where("payment_status != ?", "PAID").
		Select("COALESCE(SUM(total_amount - paid_amount), 0)").Scan(&outstandingRevenue)
	stats.OutstandingRevenue = outstandingRevenue
	
	// Growth statistics
	var newTenantsThisMonth, churnedTenantsThisMonth int64
	config.DB.Model(&models.Tenant{}).Where("created_at >= ?", firstDayOfMonth).Count(&newTenantsThisMonth)
	config.DB.Model(&models.Tenant{}).Unscoped().
		Where("deleted_at >= ?", firstDayOfMonth).Count(&churnedTenantsThisMonth)
	
	stats.NewTenantsThisMonth = int(newTenantsThisMonth)
	stats.ChurnedTenantsThisMonth = int(churnedTenantsThisMonth)
	
	// Calculate growth rate
	var lastMonthTenants int64
	firstDayOfLastMonth := firstDayOfMonth.AddDate(0, -1, 0)
	config.DB.Model(&models.Tenant{}).Where("created_at >= ? AND created_at < ?", firstDayOfLastMonth, firstDayOfMonth).Count(&lastMonthTenants)
	
	if lastMonthTenants > 0 {
		stats.GrowthRate = (float64(newTenantsThisMonth) / float64(lastMonthTenants)) * 100
	}
	
	// Usage statistics
	var totalUsers, totalCustomers int64
	config.DB.Model(&models.User{}).Count(&totalUsers)
	config.DB.Model(&models.Customer{}).Count(&totalCustomers)
	stats.TotalUsers = int(totalUsers)
	stats.TotalCustomers = int(totalCustomers)
	
	// Storage used
	var totalStorage float64
	config.DB.Model(&models.Tenant{}).Select("COALESCE(SUM(storage_used_gb), 0)").Scan(&totalStorage)
	stats.TotalStorageUsedGB = totalStorage
	
	// System statistics - defaults for now
	stats.AverageResponseTimeMs = 150.0 // TODO: Implement from metrics
	stats.ErrorRate = 0.5                // TODO: Implement from metrics
	stats.Uptime = 99.9                  // TODO: Implement from metrics
	stats.LastUpdated = time.Now()
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Platform analytics retrieved successfully",
		Data:    stats,
	})
}

// GetTenantSettings gets tenant settings
func GetTenantSettings(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	
	var settings models.TenantSettings
	if err := config.DB.Where("tenant_id = ?", tenantID).First(&settings).Error; err != nil {
		// If not found, return default settings
		if err.Error() == "record not found" {
			c.JSON(http.StatusOK, responses.SuccessResponse{
				Status:  "success",
				Message: "No custom settings found, using defaults",
				Data:    gin.H{"tenant_id": tenantID},
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch settings",
			Error:   err.Error(),
		})
		return
	}
	
	// Parse payment methods JSON
	var paymentMethods []string
	if settings.PaymentMethods != "" {
		json.Unmarshal([]byte(settings.PaymentMethods), &paymentMethods)
	}
	
	response := responses.TenantSettingsResponse{
		ID:                  settings.ID,
		TenantID:            settings.TenantID,
		CompanyName:         settings.CompanyName,
		Address:             settings.Address,
		Phone:               settings.Phone,
		Email:               settings.Email,
		Website:             settings.Website,
		LogoURL:             settings.LogoURL,
		PrimaryColor:        settings.PrimaryColor,
		SecondaryColor:      settings.SecondaryColor,
		InvoicePrefix:       settings.InvoicePrefix,
		InvoiceNumberFormat: settings.InvoiceNumberFormat,
		InvoiceDueDays:      settings.InvoiceDueDays,
		InvoiceFooterText:   settings.InvoiceFooterText,
		LatePenaltyPercent:  settings.LatePenaltyPercent,
		LatePenaltyMaxCap:   settings.LatePenaltyMaxCap,
		GracePeriodDays:     settings.GracePeriodDays,
		MinimumBillAmount:   settings.MinimumBillAmount,
		PaymentMethods:      paymentMethods,
		BankName:            settings.BankName,
		BankAccountName:     settings.BankAccountName,
		BankAccountNo:       settings.BankAccountNo,
		OperatingHours:      settings.OperatingHours,
		ServiceArea:         settings.ServiceArea,
		TimeZone:            settings.TimeZone,
		Language:            settings.Language,
		Currency:            settings.Currency,
		CreatedAt:           settings.CreatedAt,
		UpdatedAt:           settings.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Tenant settings retrieved successfully",
		Data:    response,
	})
}

// UpdateTenantSettings updates tenant settings
func UpdateTenantSettings(c *gin.Context) {
	tenantID := c.MustGet("tenant_id").(uuid.UUID)
	
	var req requests.UpdateTenantSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}
	
	var settings models.TenantSettings
	err := config.DB.Where("tenant_id = ?", tenantID).First(&settings).Error
	
	// If not found, create new settings
	if err != nil {
		settings = models.TenantSettings{
			TenantID: tenantID,
		}
	}
	
	// Update fields
	if req.CompanyName != "" {
		settings.CompanyName = req.CompanyName
	}
	if req.Address != "" {
		settings.Address = req.Address
	}
	if req.Phone != "" {
		settings.Phone = req.Phone
	}
	if req.Email != "" {
		settings.Email = req.Email
	}
	if req.Website != "" {
		settings.Website = req.Website
	}
	if req.PrimaryColor != "" {
		settings.PrimaryColor = req.PrimaryColor
	}
	if req.SecondaryColor != "" {
		settings.SecondaryColor = req.SecondaryColor
	}
	if req.InvoicePrefix != "" {
		settings.InvoicePrefix = req.InvoicePrefix
	}
	if req.InvoiceDueDays > 0 {
		settings.InvoiceDueDays = req.InvoiceDueDays
	}
	if req.InvoiceFooterText != "" {
		settings.InvoiceFooterText = req.InvoiceFooterText
	}
	if req.LatePenaltyPercent >= 0 {
		settings.LatePenaltyPercent = req.LatePenaltyPercent
	}
	if req.LatePenaltyMaxCap >= 0 {
		settings.LatePenaltyMaxCap = req.LatePenaltyMaxCap
	}
	if req.GracePeriodDays >= 0 {
		settings.GracePeriodDays = req.GracePeriodDays
	}
	if req.MinimumBillAmount >= 0 {
		settings.MinimumBillAmount = req.MinimumBillAmount
	}
	if req.BankName != "" {
		settings.BankName = req.BankName
	}
	if req.BankAccountName != "" {
		settings.BankAccountName = req.BankAccountName
	}
	if req.BankAccountNo != "" {
		settings.BankAccountNo = req.BankAccountNo
	}
	if req.OperatingHours != "" {
		settings.OperatingHours = req.OperatingHours
	}
	if req.ServiceArea != "" {
		settings.ServiceArea = req.ServiceArea
	}
	if req.TimeZone != "" {
		settings.TimeZone = req.TimeZone
	}
	if req.Language != "" {
		settings.Language = req.Language
	}
	
	if err := config.DB.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  "error",
			Message: "Failed to update settings",
			Error:   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, responses.SuccessResponse{
		Status:  "success",
		Message: "Tenant settings updated successfully",
		Data:    settings,
	})
}

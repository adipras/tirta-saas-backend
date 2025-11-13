package main

import (
	"context"
	"log"
	"os"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/middleware"
	"github.com/adipras/tirta-saas-backend/pkg/logger"
	"github.com/adipras/tirta-saas-backend/routes"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	config.Migrate()

	// Get port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	// Initialize logger
	logger.Init("INFO")
	logger.Info("🚀 Tirta-SaaS Backend starting up", map[string]interface{}{
		"version": "1.0.0",
		"env":     os.Getenv("ENV"),
		"port":    port,
	})

	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestTracingMiddleware())
	r.Use(middleware.PerformanceMonitoringMiddleware())

	// Setup Huma API for auto-documentation
	humaConfig := huma.DefaultConfig("Tirta SaaS Backend API", "1.0.0")
	humaConfig.Info.Description = `
# Multi-Tenant Water Utility Billing System

Complete SaaS solution for water utility companies (Paguyuban Air Bersih) to manage:
- **Customer Registration** - Register new water customers with meter assignment
- **Water Usage Tracking** - Record monthly meter readings and calculate usage
- **Invoice Generation** - Automated monthly billing based on usage and tariffs
- **Payment Processing** - Handle customer payments and track payment status
- **User Management** - Manage tenant admins and operators
- **Platform Management** - Platform owner tools for managing all tenants
- **Analytics & Reporting** - Real-time statistics and revenue tracking

## Authentication

Most endpoints require JWT authentication. To authenticate:

1. Register a tenant via **POST /auth/register**
2. Login via **POST /auth/login** to get JWT token
3. Include token in **Authorization** header: ` + "`Bearer <your-token>`" + `

Customer authentication:
- **POST /auth/customer/login** - Customer login endpoint

## Multi-Tenancy

Each tenant (water utility company) has complete data isolation. All operations are automatically scoped to the authenticated user's tenant.

## Rate Limiting

API endpoints are rate-limited to ensure fair usage and system stability.

## Support

For support, email: support@tirtasaas.com
`
	humaConfig.Info.Contact = &huma.Contact{
		Name:  "Tirta SaaS Support",
		Email: "support@tirtasaas.com",
		URL:   "https://tirtasaas.com/support",
	}
	humaConfig.Info.License = &huma.License{
		Name: "Proprietary",
		URL:  "https://tirtasaas.com/license",
	}
	humaConfig.Servers = []*huma.Server{
		{URL: "http://localhost:8081", Description: "Development server"},
		{URL: "https://api.tirtasaas.com", Description: "Production server"},
	}
	
	// Configure security schemes
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "JWT token from /auth/login endpoint. Include as: Authorization: Bearer <token>",
		},
	}

	// Create Huma API (auto-generates OpenAPI 3.1)
	api := humagin.New(r, humaConfig)

	// Register example Huma endpoint (with auto-documentation)
	huma.Register(api, huma.Operation{
		OperationID: "health-check-v2",
		Method:      "GET",
		Path:        "/api/v2/health",
		Summary:     "Health Check",
		Description: "Check if the API is running and healthy",
		Tags:        []string{"Health & Monitoring"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Status  string `json:"status" example:"healthy" doc:"System status"`
			Service string `json:"service" example:"tirta-saas-backend" doc:"Service name"`
			Version string `json:"version" example:"1.0.0" doc:"API version"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Status  string `json:"status" example:"healthy" doc:"System status"`
				Service string `json:"service" example:"tirta-saas-backend" doc:"Service name"`
				Version string `json:"version" example:"1.0.0" doc:"API version"`
			}
		}{}
		resp.Body.Status = "healthy"
		resp.Body.Service = "tirta-saas-backend"
		resp.Body.Version = "1.0.0"
		return resp, nil
	})

	// Legacy Swagger UI (keep for backward compatibility)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register all existing Gin routes (backward compatible)
	routes.HealthRoutes(r)
	routes.AuthRoutes(r)
	routes.ProtectedRoutes(r)
	routes.SubscriptionRoutes(r)
	routes.CustomerRoutes(r)
	routes.CustomerSelfServiceRoutes(r)
	routes.WaterRateRoutes(r)
	routes.WaterUsageRoutes(r)
	routes.InvoiceRoutes(r)
	routes.PaymentRoutes(r)
	routes.RegisterTenantUserRoutes(r)
	routes.PlatformRoutes(r)

	logger.Info("🚀 Server ready and listening", map[string]interface{}{
		"port":         port,
		"swagger":      "http://localhost:" + port + "/swagger/index.html",
		"huma_docs":    "http://localhost:" + port + "/docs",
		"openapi_json": "http://localhost:" + port + "/openapi.json",
		"openapi_yaml": "http://localhost:" + port + "/openapi.yaml",
	})
	r.Run(":" + port)
}

package main

import (
	"log"
	"os"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/middleware"
	"github.com/adipras/tirta-saas-backend/pkg/logger"
	"github.com/adipras/tirta-saas-backend/routes"
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

	// Initialize logger
	logger.Init("INFO")
	logger.Info("🚀 Tirta-SaaS Backend starting up", map[string]interface{}{
		"version": "1.0.0",
		"env":     os.Getenv("ENV"),
		"port":    port,
	})

	r := gin.Default()

	// Global middleware
	r.Use(middleware.RequestTracingMiddleware())
	r.Use(middleware.PerformanceMonitoringMiddleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check routes (must be before other middleware)
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

	logger.Info("🚀 Server ready and listening", map[string]interface{}{
		"port": port,
		"docs": "http://localhost:" + port + "/swagger/index.html",
	})
	r.Run(":" + port)
}

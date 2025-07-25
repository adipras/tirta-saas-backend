package main

import (
	"log"
	"os"

	"github.com/adipras/tirta-saas-backend/config"
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

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.AuthRoutes(r)
	routes.ProtectedRoutes(r)
	routes.SubscriptionRoutes(r)
	routes.CustomerRoutes(r)
	routes.WaterRateRoutes(r)
	routes.WaterUsageRoutes(r)
	routes.InvoiceRoutes(r)
	routes.PaymentRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}
	log.Println("🚀 Server running on port " + port)
	r.Run(":" + port)
}

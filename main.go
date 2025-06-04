package main

import (
	"log"

	"github.com/adipras/tirta-saas-backend/config"
	"github.com/adipras/tirta-saas-backend/routes"

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

	routes.AuthRoutes(r)
	routes.ProtectedRoutes(r)
	routes.SubscriptionRoutes(r)

	r.Run()
}

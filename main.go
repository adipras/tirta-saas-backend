package main

import (
	"log"

	"github.com/adipras/tirta-saas-backend/config"

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

	// Nanti kita tambahkan routes di sini
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Tirta SaaS Backend is running"})
	})

	r.Run()
}

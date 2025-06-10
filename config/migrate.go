package config

import (
	"github.com/adipras/tirta-saas-backend/models"
)

func Migrate() {
	DB.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.SubscriptionType{},
		&models.Customer{},
		&models.WaterRate{},
		&models.WaterUsage{},
		&models.Invoice{},
		&models.Payment{},
	)
}

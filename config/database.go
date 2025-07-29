package config

import (
	"fmt"
	"log"
	"os"

	"github.com/adipras/tirta-saas-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || port == "" || name == "" {
		log.Fatal("❌ ENV database tidak lengkap. Harap periksa .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("❌ Gagal konek ke database: %v", err)
	}

	log.Println("✅ Database berhasil terkoneksi")
	DB = database
}

func Migrate() {
	log.Println("🚀 Memulai proses migrasi database...")

	// Migration order is important due to foreign key constraints
	// 1. Base entities first (no dependencies)
	// 2. Entities with foreign keys last
	err := DB.AutoMigrate(
		&models.Tenant{},            // No dependencies
		&models.User{},              // References Tenant
		&models.SubscriptionType{},  // References Tenant
		&models.Customer{},          // References Tenant + SubscriptionType
		&models.WaterRate{},         // References Tenant + SubscriptionType
		&models.WaterUsage{},        // References Tenant + Customer
		&models.Invoice{},           // References Tenant + Customer
		&models.Payment{},           // References Tenant + Invoice (must be after Invoice)
		&models.AuditLog{},          // References Tenant (no other FK constraints)
	)

	if err != nil {
		log.Fatalf("❌ Migrasi gagal: %v", err)
	}

	log.Println("✅ Migrasi database selesai.")
	
	// Apply database optimizations after migration
	if err := OptimizeDatabase(DB); err != nil {
		log.Printf("⚠️ Database optimization failed: %v", err)
	} else {
		log.Println("✅ Database optimizations applied")
	}
}

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

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal konek ke database: %v", err)
	}

	log.Println("✅ Database berhasil terkoneksi")
	DB = database
}

func Migrate() {
	log.Println("🚀 Memulai proses migrasi database...")

	err := DB.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.SubscriptionType{},
		&models.Customer{},
		&models.WaterRate{},
		&models.WaterUsage{},
		&models.Invoice{},
		&models.Payment{},
	)

	if err != nil {
		log.Fatalf("❌ Migrasi gagal: %v", err)
	}

	log.Println("✅ Migrasi database selesai.")
}

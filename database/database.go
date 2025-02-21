package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"taskmanager/models"
)

var DB *gorm.DB

func ConnectDB() {
	// (Optionally) Load .env (if not already loaded in main.go)
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env file already loaded or not found, continuing")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL is not set in environment variables")
	}

	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	fmt.Println("✅ Successfully connected to the database")

	// AutoMigrate the models (creates/updates tables)
	err = DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal("❌ Failed to run migrations:", err)
	}
	fmt.Println("✅ Database Migration Completed")
}

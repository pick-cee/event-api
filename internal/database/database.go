package database

import (
	"fmt"
	"log"

	"github.com/pick-cee/events-api/internal/config"
	"github.com/pick-cee/events-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	return nil
}

func Migrate() error {
	log.Println("🔄 Running migrations...")
  
	err := DB.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Registration{},
	)

	if err != nil {
		return err
	}

	log.Println("✅ Migrations completed")
  return nil
}

func Disconnect() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	log.Println("🔌 Disconnecting from database...")
	return sqlDB.Close()
}
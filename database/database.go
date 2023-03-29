package database

import (
	"discord-bot-dashboard-backend-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Dsn string
}

func Start(config Config) *gorm.DB {
	dsn := config.Dsn
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&models.Guild{}); err != nil {
		panic("failed to migrate models")
	}

	return db
}

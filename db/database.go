package db

import (
	"discord-bot-dashboard-backend-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataBaseConfig struct {
	Host     string
	Name     string
	User     string
	Password string
}

func Start(config DataBaseConfig) *gorm.DB {
	dsn := "host=" + config.Host + " user=" + config.User + " password=" + config.Password + " dbname=" + config.Name
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&models.Guild{}); err != nil {
		panic("failed to migrate models")
	}

	return db
}

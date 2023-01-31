package database

import (
	"discord-bot-dashboard-backend-go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
	Dsn      *string
}

func Start(config Config) *gorm.DB {
	var dsn string

	if config.Dsn == nil {
		dsn = "host=" + config.Host + " user=" + config.User + " password=" + config.Password + " dbname=" + config.Name + " port=" + config.Port + " sslmode=require"
	} else {
		dsn = *config.Dsn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&models.Guild{}); err != nil {
		panic("failed to migrate models")
	}

	return db
}

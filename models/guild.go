package models

type Guild struct {
	Id             string  `gorm:"primaryKey"`
	WelcomeChannel *string `gorm:"default: null"`
	WelcomeMessage *string `gorm:"default: null"`
}

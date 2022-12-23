package models

type Guild struct {
	Id             string  `gorm:"primaryKey"`
	WelcomeMessage *string `gorm:"default: null"`
}

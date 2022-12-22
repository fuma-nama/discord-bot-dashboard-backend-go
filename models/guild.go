package models

import (
	"gorm.io/gorm"
)

type Guild struct {
	gorm.Model
	Id uint `gorm:"primaryKey"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type RecievedEmail struct {
	gorm.Model
	FromEmail string `gorm:"not null"`
	ToEmail string `gorm:"not null"`
	RecievedAt time.Time `gorm:"not null"`
	Contnet string `gorm:"not null"`
	Amount int64 `gorm:"not null"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type DailyTotal struct {
	gorm.Model
	UserEmail string `gorm:"not null;index"` 
	Date time.Time `gorm:"not null;index"`
	Amount int64 `gorm:"not null"`
}

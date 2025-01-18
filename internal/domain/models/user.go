package models 

import (
	"time"

	"gorm.io/gorm"
)

// データベースの構造に依存せず、ビジネスロジックに必要な形で定義します
type User struct {
	gorm.Model
	Email string `gorm:"uniqueIndex;not nul"`
	AccessToken   string `gorm:"not null"`
	RefreshToken  string `gorm:"not null"`
	TokenExpiry   time.Time
}




package handlers

import "gorm.io/gorm"


type HealthHandler struct {
	db *gorm.DB
}

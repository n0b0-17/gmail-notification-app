package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context){
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"message": "Failed to get batabase instance",
		})
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"message": "Database connection lost",
		})
		return
	}

	c.JSON(200,gin.H{
		"status":"ok",
		"database": "connected",
	})
}


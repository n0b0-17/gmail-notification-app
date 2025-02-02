package handlers

import (
	"fmt"
	"gmail-notification-app/internal/domain/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DailyTotalHandler struct {
	db *gorm.DB
}


func NewDailyTotalHandler(db *gorm.DB) *DailyTotalHandler {
	return &DailyTotalHandler{db: db}
}

func (h *DailyTotalHandler) GetDailyTotalHandler(c *gin.Context) {
	session := sessions.Default(c)
	userEmail := session.Get("user_email").(string)

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	today := time.Now().In(jst).Truncate(24 * time.Hour)

	var totalAmount int64
	result := h.db.Model(&models.RecievedEmail{}).
		Where("to_email = ? AND DATE(recieved_at) = DATE(?)", userEmail, today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate daily total"})
		return
	}

	dailyTotal := &models.DailyTotal{
		UserEmail: userEmail,
		Date: today,
		Amount: totalAmount,
	}

	if err := h.db.Create(dailyTotal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save daily total"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date": today.Format("01/02"),
		"amount": totalAmount,
		"message": fmt.Sprintf("%s日は%d円使用しました", today.Format("01/02"), totalAmount),
	})
}


func (h *DailyTotalHandler) BatchGetDailyTotal(c *gin.Context) {
    apiKey := c.GetHeader("X-API-Key")
    if apiKey != os.Getenv("BATCH_API_KEY") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
        return
    }

    // ユーザー全員分の集計を行う
    var users []models.User
    if err := h.db.Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }

    jst := time.FixedZone("Asia/Tokyo", 9*60*60)
    today := time.Now().In(jst).Truncate(24 * time.Hour)

    for _, user := range users {
        var totalAmount int64
        result := h.db.Model(&models.RecievedEmail{}).
            Where("to_email = ? AND DATE(recieved_at) = DATE(?)", user.Email, today).
            Select("COALESCE(SUM(amount), 0)").
            Scan(&totalAmount)

        if result.Error != nil {
            log.Printf("Failed to calculate daily total for user %s: %v", user.Email, result.Error)
            continue
        }

        dailyTotal := &models.DailyTotal{
            UserEmail: user.Email,
            Date:      today,
            Amount:    totalAmount,
        }

        if err := h.db.Create(dailyTotal).Error; err != nil {
            log.Printf("Failed to save daily total for user %s: %v", user.Email, err)
        }
    }

    c.JSON(http.StatusOK, gin.H{"message": "Daily totals calculated successfully"})
}

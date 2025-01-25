package handlers

import (
	"gmail-notification-app/internal/domain/models"
	"gmail-notification-app/internal/infrastructure/gmail"
	"gmail-notification-app/internal/infrastructure/oauth"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GmailHandler struct{
	googleOAuth *oauth.GoogleOAuth
	db *gorm.DB
}

func NewGmailHandler(googleOauth *oauth.GoogleOAuth,db *gorm.DB) *GmailHandler {
	return &GmailHandler{
		googleOAuth: googleOauth,
		db: db,
	}
}

func (h *GmailHandler) ListEmails(c *gin.Context) {
	var user models.User

	if err := h.db.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})	
		return
	}

	service, err := gmail.NewGmailService(user.AccessToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Gmail Serivice"})
		return
	}
	messages,err := service.ListMessages("post_master@netbk.co.jp")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	var messageDetails []map[string]interface{}
	for _, msg := range messages{
		headers := msg.Payload.Headers
		detail := make(map[string]interface{})

		for _, header := range headers{
			switch header.Name {
			case "From":
				detail["from"] = header.Value
			case "Subject":
				detail["subject"] = header.Value
			case "Date":
				detail["date"] = header.Value
			}			
			detail["id"] = msg.Id
			messageDetails = append(messageDetails, detail)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messageDetails,
	})

}

package handlers

import (
	"gmail-notification-app/internal/domain/models"
	"gmail-notification-app/internal/infrastructure/gmail"
	"gmail-notification-app/internal/infrastructure/oauth"
	"net/http"

	"github.com/gin-contrib/sessions"
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

func (h *GmailHandler) SavedEmails(c *gin.Context) {
	sessions := sessions.Default(c)
	userEmail := sessions.Get("user_email").(string)
	var user models.User

	if err := h.db.Where("email = ?",userEmail).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})	
		return
	}

	service, err := gmail.NewGmailService(user.AccessToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Gmail Serivice"})
		return
	}
	messages,err := service.ListMessages("test.nob.dev@gmail.com")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}


	//DBに保存するトランザクション
	tx := h.db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to begin transaction"})
		return 
	}

	var savedEmails []models.RecievedEmail

	for _, msg  := range messages{
		email := gmail.ConvertToRecievedEmail(msg)

		var existingEmail models.RecievedEmail
		if err := h.db.Where("recieved_at = ? AND to_email = ?",email.RecievedAt,userEmail).First(&existingEmail).Error; err == nil{
				continue
			}	

		if err := tx.Create(email).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to save emails"})
			return
		}
		savedEmails = append(savedEmails, *email)
	}

	//トランザクションのコミット
	if err := tx.Commit().Error; err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Emails saved successfully",
		"count": len(savedEmails),
	})

}

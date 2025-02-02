package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/gin-contrib/sessions"
)


func AuthMiddleware(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		session := sessions.Default(c)
		userEmail := session.Get("user_email")
		if userEmail == nil{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User email not provided"})
			c.Abort()
			return
		}

		c.Set("user_email", userEmail)
		c.Next()

	}
}

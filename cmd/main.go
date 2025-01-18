package main

import (
	"gmail-notification-app/config"
	"gmail-notification-app/internal/infrastructure/database"
	"gmail-notification-app/internal/infrastructure/oauth"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	//データベース接続
	db, err := database.NewDatabase(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//Google OAuth設定
	GoogleOAuth := oauth.NewGoogleOAuth(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.RedirectURL,
	)

  r := gin.Default()
    
    // ヘルスチェックエンドポイント
    r.GET("/health", func(c *gin.Context) {
			sqlDB, err := db.DB()
			if err != nil {
				c.JSON(500,gin.H{
					"status":"error",
					"message": "Failed to get database instance",
				})
				return
			}
			err = sqlDB.Ping()
			if err != nil {
				c.JSON(500,gin.H{
					"status":"error",
					"message": "Database connection lost",
				})
				return
			}
			c.JSON(200,gin.H{
				"status":"ok",
				"message":"connected"
			})
    })

		r.GET("/auth/login", func(c *gin.Context) {
			authURL := GoogleOAuth.GetAuthURL()
			c.Redirect(302,authURL)
		})

		log.Println("Starting server on Port:8080")
		r.Run(":8080")
}


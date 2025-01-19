package main

import (
	"gmail-notification-app/config"
	"gmail-notification-app/internal/infrastructure/database"
	"gmail-notification-app/internal/infrastructure/oauth"
	"gmail-notification-app/internal/interfaces/handlers"
	"gmail-notification-app/internal/interfaces/router"
	"log"

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

	if GoogleOAuth == nil {
		log.Fatalf("Failed to initialized GoogleOAuth")
	}

	//ハンドラーの初期化
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(GoogleOAuth,db)

	//ルーターのセットアップ
	r := router.SetupRouter(healthHandler,authHandler)

	log.Println("Starting server on Port:8080")
	if err := r.Run(":8080"); err != nil{
		log.Fatalf("Server failed to start: %v", err)
	}
}


package main

import (
	"gmail-notification-app/config"
	"gmail-notification-app/internal/infrastructure/bot"
	"gmail-notification-app/internal/infrastructure/database"
	"gmail-notification-app/internal/infrastructure/oauth"
	"gmail-notification-app/internal/infrastructure/scheduler"
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


	//Line Messaging APIの初期化

	lineClient, err := bot.NewLineBot(cfg.LineSecretID, cfg.LineToken)

	if err != nil {
		log.Fatalf("Failed to initialized LineBot")
	}


  routerCfg := &router.RouterConfig{
		DB: db,
		HealthHandler: handlers.NewHealthHandler(db),
		AuthHandler: handlers.NewAuthHandler(GoogleOAuth, db),
		GmailHandler: handlers.NewGmailHandler(GoogleOAuth, db),
		DailyTotalHandler: handlers.NewDailyTotalHandler(db),
		NotifyHandler: handlers.NewNotifyHandler(lineClient),
  }

	//ルーターのセットアップ
	r := router.SetupRouter(*routerCfg)
	scheduler := scheduler.NewScheduler("http://localhost:8080")
	
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer scheduler.Stop()

	if err := r.Run(":8080"); err != nil{
		log.Fatalf("Server failed to start: %v", err)
	}

}


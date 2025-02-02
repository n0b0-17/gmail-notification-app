package router

import (
	"gmail-notification-app/internal/interfaces/handlers"
	"gmail-notification-app/internal/interfaces/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RouterConfig struct{
	DB *gorm.DB
	HealthHandler *handlers.HealthHandler
	AuthHandler *handlers.AuthHandler
	GmailHandler *handlers.GmailHandler
	DailyTotalHandler *handlers.DailyTotalHandler
	NotifyHandler *handlers.NotifyHandler
}


func SetupRouter(cfg RouterConfig) *gin.Engine{
	r := gin.Default()
	
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/health",cfg.HealthHandler.Check)

	auth := r.Group("/auth")
	{
		auth.GET("/login",cfg.AuthHandler.Login)
		auth.GET("/callback",cfg.AuthHandler.Callback)
	}

	gmail := r.Group("/gmail")
	gmail.Use(middleware.AuthMiddleware(cfg.DB))
	{
		gmail.GET("/saved", cfg.GmailHandler.SavedEmails)
		gmail.GET("daily-totals",cfg.DailyTotalHandler.GetDailyTotalHandler)
	}

	notify := r.Group("notify")
	{
		notify.GET("/test",cfg.NotifyHandler.SendMessage)
	}

	batch := r.Group("/batch")
	{
		batch.GET("/daily-totals", cfg.DailyTotalHandler.BatchGetDailyTotal)
	}
	return r
}

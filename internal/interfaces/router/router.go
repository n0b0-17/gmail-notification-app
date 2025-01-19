package router

import (
	"gmail-notification-app/internal/interfaces/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler) *gin.Engine{
	r := gin.Default()
	
	// パニックリカバリーミドルウェアを追加
	r.Use(gin.Recovery())

	r.GET("/health",healthHandler.Check)

	auth := r.Group("/auth")
	{
		auth.GET("/login",authHandler.Login)
		auth.GET("/callback",authHandler.Callback)
	}
	return r
}

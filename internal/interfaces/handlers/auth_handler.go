package handlers

import (
	"gmail-notification-app/internal/infrastructure/oauth"

	"github.com/gin-gonic/gin"
)


type AuthHandler struct{
 googleOAuth	*oauth.GoogleOAuth
}

func NewAuthHandler(googleOAuth *oauth.GoogleOAuth) *AuthHandler {
	/*
	~~~~~~~~~~~~~~~~~~
	Role:
		認証関連のハンドラーを生成するコンストラクタ
	~~~~~~~~~~~~~~~~~~
	*/
    return &AuthHandler{googleOAuth: googleOAuth}
}

func (h *AuthHandler)Login(c *gin.Context){
	authURL := h.googleOAuth.GetAuthURL()
	c.Redirect(302,authURL)
}

func (h *AuthHandler) Callback(c *gin.Context) {
    c.JSON(200, gin.H{"message": "callback endpoint"})
}


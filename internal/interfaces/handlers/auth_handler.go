package handlers

import (
	"fmt"
	"gmail-notification-app/internal/domain/models"
	"gmail-notification-app/internal/infrastructure/oauth"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)


type AuthHandler struct{
 googleOAuth	*oauth.GoogleOAuth
 db *gorm.DB
}

func NewAuthHandler(googleOAuth *oauth.GoogleOAuth,db *gorm.DB) *AuthHandler {
	/*
	~~~~~~~~~~~~~~~~~~
	Role:
		認証関連のハンドラーを生成するコンストラクタ
	~~~~~~~~~~~~~~~~~~
	*/
	return &AuthHandler{
		googleOAuth: googleOAuth,
		db: db,
	}
}

func (h *AuthHandler)Login(c *gin.Context){
	// パニックをリカバー
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Login handler: %v", r)
			debug.PrintStack()
			c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error occurred",
			})
		}
	}()
	if h.googleOAuth == nil {
		log.Printf("GoogleOAuth is nil in Login handler")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "OAuth configuration is missing",
		})
		return
	}
	authURL := h.googleOAuth.GetAuthURL()
	log.Printf("Generated auth URL: %s", authURL)


	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization code is missing",
		})
		return
	}

	token,err := h.googleOAuth.ExchangeCodeForToken(code)
	
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to exchange token: %v", err),
				"code_length": len(code),  
		})
		return
	}

	userInfo, err := h.googleOAuth.FethcedUserInfo(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errror": "Failed to fetch user info",
		})
		return
	}

	//データベースに保存
	if err := h.saveUser(userInfo.Email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save user information",
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token recieved",
		"expires_at": userInfo.Email,
	})
}

func (h *AuthHandler) saveUser(email string,token *oauth2.Token) error {
	user := &models.User{
		Email: email,
		AccessToken: token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry: token.Expiry,
	}
	
	return h.db.Transaction(func(tx *gorm.DB) error {
		var existingUser models.User

		result := tx.Where("email = ?", email).First(&existingUser)

		if result.Error == gorm.ErrRecordNotFound {
			//新規ユーザーの場合は作成
			return tx.Create(user).Error
		}else if result.Error != nil{
			return result.Error
		}

		//既存ユーザーの場合はトークンを更新
		updates := map[string]interface{}{
			"access_token": token.AccessToken,
			"refresh_token": token.RefreshToken,
			"token_expiry": token.Expiry,
		}
		return tx.Model(&existingUser).Updates(updates).Error
	})
}

func (h *AuthHandler) UpdateUserToken(email string) error {
	var user models.User
	if err := h.db.Where("email =?",email).First(&user).Error; err != nil {
		return fmt.Errorf("use not found: %w", err)
	}

	currentToken := oauth.TokenInfo{
		AccessToken: user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry: user.TokenExpiry,
	}

	//トークンの更新
	newToken, err := h.googleOAuth.RefreshAccessToken(currentToken)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	//データベースの更新
	updates := map[string]interface{}{
		"access_token": newToken.AccessToken,
		"refresh_token": newToken.RefreshToken,
		"token_expiry": newToken.Expiry,
	}
	
	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user token: %w", err)
	}
	return nil
}


package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	config *oauth2.Config
}

type UserInfo struct {
    Email string `json:"email"`
    Id    string `json:"id"`
} 
type TokenInfo struct {
	AccessToken string
	RefreshToken string
	Expiry time.Time
}

func NewGoogleOAuth(clientID, clientSecret,redirectURL string) *GoogleOAuth {
	// パラメータの検証
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		log.Printf("Invalid OAuth parameters - ClientID empty: %v, ClientSecret empty: %v, RedirectURL empty: %v",
			clientID == "",
			clientSecret == "",
			redirectURL == "",
		)
		return nil
	}
	return &GoogleOAuth{
		config: &oauth2.Config{
			ClientID: clientID,
			ClientSecret: clientSecret,
			RedirectURL: redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/gmail.readonly",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *GoogleOAuth) GetAuthURL() string{
	return g.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (g *GoogleOAuth) Config() *oauth2.Config{
	return g.config
}


func (g *GoogleOAuth) ExchangeCodeForToken(code string) (*oauth2.Token, error){
	token, err := g.config.Exchange(oauth2.NoContext, code)
	if err != nil {
			return nil, fmt.Errorf("oauth exchange error: %w", err)
	}
	return token, nil
}

func (g *GoogleOAuth) FethcedUserInfo(accessToken string) (*UserInfo, error){
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("failed to fetch user info :%s", resp.Status)
	}

	var userInfo UserInfo

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo,nil
}

func (g *GoogleOAuth) RefreshAccessToken(currentToken TokenInfo) (*TokenInfo,error){
	token := &oauth2.Token{
		AccessToken: currentToken.AccessToken,
		RefreshToken: currentToken.RefreshToken,
		Expiry: currentToken.Expiry,
    TokenType:    "Bearer",
	}

	//Tokenの有効期限までに60分の猶予がある場合は更新をしない
	if token.Expiry.Add(-60 * time.Minute).After(time.Now()){
		return &currentToken,nil 
	}

	//トークンを更新
	newToken, err := g.config.TokenSource(context.Background(),token).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w",err)
	}

	refreshToken := newToken.RefreshToken
	//Memo:既存のトークンをTokenSource().Token()で更新かけた場合
	//     基本的にはあたらしいrefreshTokenを取得しないため、空になる可能性がある
	//     その場合は、既存のリフレッシュトークンを引き続き使用する
	if refreshToken == "" {
		refreshToken = currentToken.RefreshToken
	}

	return &TokenInfo{
		AccessToken: newToken.AccessToken,
		RefreshToken: refreshToken,
		Expiry: newToken.Expiry,
	},nil
}





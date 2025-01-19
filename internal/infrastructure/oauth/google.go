package oauth

import (
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	config *oauth2.Config
}

func NewGoogleOAuth(clientID, clientSecret,redirectURL string) *GoogleOAuth {
	log.Printf("Initializing Google OAuth with:")
  log.Printf("Client ID: %s", clientID[:5]+"...") // セキュリティのため全文は表示しない
  log.Printf("Redirect URL: %s", redirectURL)
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



package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct{
	GoogleClientID string
	GoogleClientSecret string
	RedirectURL string
	DatabaseURL string
    LineSecretID string 
    LineToken string
}

func init() {
    // .envファイルの読み込み
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found, using system environment variables")
    }
}

func NewConfig() *Config {

    cfg := &Config{
        GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:        os.Getenv("REDIRECT_URL"),
        DatabaseURL:        os.Getenv("DATABASE_URL"),
        LineSecretID: os.Getenv("CHANNEL_SECRET"),
        LineToken: os.Getenv("CHANNEL_TOKEN"),
    }
    
    // デバッグ用のログ追加
    log.Printf("Configuration loaded:")
    log.Printf("- Client ID exists: %v", cfg.GoogleClientID != "")
    log.Printf("- Client Secret exists: %v", cfg.GoogleClientSecret != "")
    log.Printf("- Redirect URL: %s", cfg.RedirectURL)
    log.Printf("- LINE Channel ID exists: %v", cfg.LineSecretID != "")
    log.Printf("- LINE Channel Token exists: %v", cfg.LineToken != "")


    if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" {
        log.Fatal("Google OAuth credentials are missing")
    }
    if cfg.RedirectURL == "" {
        log.Fatal("Redirect URL is missing")
    }

    return cfg
}



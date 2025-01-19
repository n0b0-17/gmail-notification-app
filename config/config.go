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
    }
    // デバッグ用のログ追加
    log.Printf("Configuration loaded:")
    log.Printf("- Client ID exists: %v", cfg.GoogleClientID != "")
    log.Printf("- Client Secret exists: %v", cfg.GoogleClientSecret != "")
    log.Printf("- Redirect URL: %s", cfg.RedirectURL)

    // 値の検証を追加
    if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" {
        log.Fatal("Google OAuth credentials are missing")
    }
    if cfg.RedirectURL == "" {
        log.Fatal("Redirect URL is missing")
    }

    return cfg
}



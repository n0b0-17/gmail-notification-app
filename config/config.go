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

    // 必須値の検証
    if cfg.GoogleClientID == "" {
        log.Fatal("GOOGLE_CLIENT_ID is not set")
    }
    if cfg.GoogleClientSecret == "" {
        log.Fatal("GOOGLE_CLIENT_SECRET is not set")
    }
    if cfg.RedirectURL == "" {
        log.Fatal("REDIRECT_URL is not set")
    }

    return cfg
}



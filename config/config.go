package config

import "os"

type Config struct{
	GoogleClientID string
	GoogleClientSecret string
	RedirectURL string
	DatabaseURL string
}


func NewConfig() *Config {
	return &Config{
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECLET"),
		RedirectURL: os.Getenv("REDIRECT_URL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}



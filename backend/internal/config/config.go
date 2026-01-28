package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string

	// Meta API
	MetaAccessToken    string
	MetaAppSecret      string
	MetaVerifyToken    string
	WhatsAppPhoneID    string
	WhatsAppBusinessID string
	InstagramAccountID string

	// n8n Integration
	N8NWebhookURL string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/omnichannel"),
		Port:        getEnv("PORT", "8080"),

		MetaAccessToken:    os.Getenv("META_ACCESS_TOKEN"),
		MetaAppSecret:      os.Getenv("META_APP_SECRET"),
		MetaVerifyToken:    getEnv("META_VERIFY_TOKEN", "omnichannel_verify_token"),
		WhatsAppPhoneID:    os.Getenv("WHATSAPP_PHONE_ID"),
		WhatsAppBusinessID: os.Getenv("WHATSAPP_BUSINESS_ID"),
		InstagramAccountID: os.Getenv("INSTAGRAM_ACCOUNT_ID"),

		N8NWebhookURL: os.Getenv("N8N_WEBHOOK_URL"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

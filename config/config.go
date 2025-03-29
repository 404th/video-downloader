package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	BotMode          string `json:"bot_mode"`
	TelegramBotToken string `json:"telegram_bot_token"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = new(Config)

	godotenv.Load("./.env")

	cfg.BotMode = cast.ToString(getEnvOrSetDefault("BOT_MODE", "debug"))
	cfg.TelegramBotToken = cast.ToString(getEnvOrSetDefault("TELEGRAM_BOT_TOKEN", "secret_telegram_bot_token"))

	return
}

func getEnvOrSetDefault(key string, defaultValue any) any {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

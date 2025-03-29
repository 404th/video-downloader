package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	BotMode                string `json:"bot_mode"`
	TelegramBotToken       string `json:"telegram_bot_token"`
	TelegramUsernameChatId int64  `json:"telegram_username_chat_id"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = new(Config)

	godotenv.Load("./.env")

	cfg.BotMode = cast.ToString(getEnvOrSetDefault("BOT_MODE", "debug"))
	cfg.TelegramBotToken = cast.ToString(getEnvOrSetDefault("TELEGRAM_BOT_TOKEN", "secret_telegram_bot_token"))
	cfg.TelegramUsernameChatId = cast.ToInt64(getEnvOrSetDefault("TELEGRAM_USERNAME_CHAT_ID", 12323452))

	return
}

func getEnvOrSetDefault(key string, defaultValue any) any {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

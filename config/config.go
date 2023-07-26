package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	TelegramToken      string
	TelegramWebHookURL string
	DbUrl              string
}

func GetConfig() EnvConfig {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env getting from actual env")
	}

	return EnvConfig{
		TelegramToken:      os.Getenv("TELEGRAM_TOKEN"),
		TelegramWebHookURL: os.Getenv("TELEGRAM_WEBHOOK_URL"),
		DbUrl:              os.Getenv("DB_URL"),
	}
}

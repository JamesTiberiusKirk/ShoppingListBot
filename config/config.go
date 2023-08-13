package config

import (
	"os"
	"strconv"

	log "github.com/inconshreveable/log15"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	TelegramToken      string
	TelegramWebHookURL string
	DbUrl              string
	Debug              bool
}

func GetConfig() EnvConfig {
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env getting from actual env")
	}

	return EnvConfig{
		TelegramToken:      os.Getenv("TELEGRAM_TOKEN"),
		TelegramWebHookURL: os.Getenv("TELEGRAM_WEBHOOK_URL"),
		DbUrl:              os.Getenv("DB_URL"),
		Debug:              getenvBool("DEBUG"),
	}
}

func getenvBool(key string) bool {
	s := os.Getenv(key)
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}

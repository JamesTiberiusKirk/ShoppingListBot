package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	TelegramToken string
	DbUrl         string
	// DbName        string
	// DbHost        string
	// DbUser        string
	// DbPass        string
}

func GetConfig() EnvConfig {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env getting from actual env")
	}

	return EnvConfig{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DbUrl:         os.Getenv("DB_URL"),
		// DbName:        os.Getenv("DB_NAME"),
		// DbHost:        os.Getenv("DB_HOST"),
		// DbUser:        os.Getenv("DB_USER"),
		// DbPass:        os.Getenv("DB_PASS"),
	}
}

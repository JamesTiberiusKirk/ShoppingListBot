package main

import (
	"log"
	"os"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/bot"
	"github.com/joho/godotenv"
)

type config struct {
	TelegramToken string
	DbHost        string
	DbUser        string
	DbPass        string
}

func getConfig() config {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env getting from actual env")
	}

	return config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DbHost:        os.Getenv("DB_HOST"),
		DbUser:        os.Getenv("DB_USER"),
		DbPass:        os.Getenv("DB_PASS"),
	}
}

func main() {
	c := getConfig()

	err := bot.StartBot(c.TelegramToken, true)
	if err != nil {
		panic(err)
	}
}

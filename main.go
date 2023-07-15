package main

import (
	"log"
	"os"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/bot"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/clients"
	"github.com/joho/godotenv"
)

type config struct {
	TelegramToken string
	DbName        string
	DbHost        string
	DbUser        string
	DbPass        string
	DbUrl         string
}

func getConfig() config {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env getting from actual env")
	}

	return config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DbName:        os.Getenv("DB_NAME"),
		DbHost:        os.Getenv("DB_HOST"),
		DbUser:        os.Getenv("DB_USER"),
		DbPass:        os.Getenv("DB_PASS"),
		DbUrl:         os.Getenv("DB_URL"),
	}
}

func main() {
	c := getConfig()

	db, err := clients.NewDBClient(c.DbUrl)
	if err != nil {
		panic(err)
	}

	err = bot.StartBot(c.TelegramToken, true, db)
	if err != nil {
		panic(err)
	}
}

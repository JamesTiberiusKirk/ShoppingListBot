package main

import (
	"github.com/JamesTiberiusKirk/ShoppingListsBot/bot"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/config"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/db"
)

func main() {
	c := config.GetConfig()

	db, err := db.NewDBClient(c.DbUrl)
	if err != nil {
		panic(err)
	}

	err = bot.StartBot(c.TelegramToken, c.TelegramWebHookURL, true, db)
	if err != nil {
		panic(err)
	}
}

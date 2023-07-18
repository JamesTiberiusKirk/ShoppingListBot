package main

import (
	"github.com/JamesTiberiusKirk/ShoppingListsBot/bot"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/clients"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/config"
)

func main() {
	c := config.GetConfig()

	db, err := clients.NewDBClient(c.DbUrl)
	if err != nil {
		panic(err)
	}

	err = bot.StartBot(c.TelegramToken, true, db)
	if err != nil {
		panic(err)
	}
}

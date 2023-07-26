package main

import (
	"github.com/JamesTiberiusKirk/ShoppingListsBot/bot"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/config"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/db"
	log "github.com/inconshreveable/log15"
)

func main() {
	c := config.GetConfig()

	db, err := db.NewDBClient(c.DbUrl)
	if err != nil {
		panic(err)
	}

	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)

	err = bot.StartBot(c.TelegramToken, c.TelegramWebHookURL, false, db)
	if err != nil {
		panic(err)
	}
}

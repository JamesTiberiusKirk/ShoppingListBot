package main

import (
	"github.com/JamesTiberiusKirk/ShoppingListsBot/config"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/db"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/handlers"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/tgf"
)

var version = "development"

func GetVersion() string {
	return version
}

func main() {
	c := config.GetConfig()

	db, err := db.NewDBClient(c.DbUrl)
	if err != nil {
		panic(err)
	}

	js := handlers.NewDBJourneyStore(db)

	botAPI, err := tgf.InitBotAPI(c.TelegramToken, c.TelegramWebHookURL, false)
	if err != nil {
		panic(err)
	}

	commands := handlers.GetHandlerCommandList()
	journeys := handlers.NewHandlerJounreyMap(botAPI, db, GetVersion)

	bot := tgf.NewBot(botAPI, commands, journeys, nil, js)
	err = bot.StartBot(false, db)
	if err != nil {
		panic(err)
	}
}

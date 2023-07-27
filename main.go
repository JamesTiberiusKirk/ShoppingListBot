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

	// h := log.CallerFileHandler(log.StdoutHandler)
	// log.Root().SetHandler(h)

	logger := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(logger)

	// logger = log.LvlFilterHandler(
	// 	log.LvlError,
	// 	log.StreamHandler(os.Stderr, log.TerminalFormat()),
	// )
	// log.Root().SetHandler(logger)
	//
	// logger = log.LvlFilterHandler(
	// 	log.LvlInfo,
	// 	log.StreamHandler(os.Stdout, log.TerminalFormat()),
	// )
	// log.Root().SetHandler(logger)

	// errorHandler := log.LvlFilterHandler(log.LvlError)
	// infoHandler := log.LvlFilterHandler(log.LvlInfo, h)
	// log.Root().SetHandler(h)

	err = bot.StartBot(c.TelegramToken, c.TelegramWebHookURL, false, db)
	if err != nil {
		panic(err)
	}
}

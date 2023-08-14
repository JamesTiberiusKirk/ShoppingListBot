package handlers

import (
	"github.com/JamesTiberiusKirk/ShoppingListsBot/db"
	"github.com/JamesTiberiusKirk/tgf"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetHandlerCommandList() tgbotapi.SetMyCommandsConfig {
	return tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "version",
			Description: "Print version of the bot",
		},
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Welcome and chat registration",
		},
		tgbotapi.BotCommand{
			Command:     "newlist",
			Description: "Create new shopping list",
		},
		tgbotapi.BotCommand{
			Command:     "additems",
			Description: "Add items to a shopping list",
		},
		tgbotapi.BotCommand{
			Command:     "displaylist",
			Description: "Display shopping list",
		},
		tgbotapi.BotCommand{
			Command:     "editlists",
			Description: "Display and edit lists",
		},
	)
}

func NewHandlerJounreyMap(bot *tgbotapi.BotAPI, db *db.DB, getVer func() string) map[string]tgf.HandlerInterface {
	return map[string]tgf.HandlerInterface{
		"version":  NewVersionHandler(bot.Send, getVer),
		"start":    NewStartHandler(bot.Send, db.AddNewChat, db.CheckIfChatExists),
		"newlist":  NewNewListHandler(bot.Send, db.NewShoppingList, db.CheckIfChatExists),
		"additems": NewAddItemsHandler(bot.Send, bot.Request, db.GetListsByChat, db.AddItemsToList, db.CheckIfChatExists),
		"displaylist": NewDisplayListHandler(bot.Send, bot.Request, db.GetListsByChat,
			db.GetItemsByList, db.ToggleItemPurchase, db.CheckIfChatExists, db.DeleteItem),
		"editlists": NewEditListsHandler(bot.Send, bot.Request, db.GetListsByChat,
			db.GetItemsByList, db.CheckIfChatExists, db.DeleteListByID),
	}
}

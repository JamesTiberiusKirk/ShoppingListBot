package handlers

import (
	"github.com/JamesTiberiusKirk/ShoppingListsBot/clients"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handlerFunc func(update tgbotapi.Update) error

func GetHandlerCommandList() tgbotapi.SetMyCommandsConfig {
	return tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Welcome and chat registration",
		},
		tgbotapi.BotCommand{
			Command:     "kb",
			Description: "Get keyboard",
		},
		tgbotapi.BotCommand{
			Command:     "newlist",
			Description: "Create new shopping list",
		},
	)
}

func NewHandlerMap(bot *tgbotapi.BotAPI, db *clients.DB) map[string]handlerFunc {
	return map[string]handlerFunc{
		"start":   NewStartHandler(bot, db.AddNewChat, db.CheckIfChatExists).Handle,
		"newlist": NewNewListHandler(bot).Handle,
		"kb":      NewKeyboardHandler(bot).Handle,
	}
}

func NewCallbackMap(bot *tgbotapi.BotAPI, db *clients.DB) map[string]handlerFunc {
	return map[string]handlerFunc{
		"kb": NewKeyboardHandler(bot).Callback,
	}
}

func NewReplyCallbackMap(bot *tgbotapi.BotAPI, db *clients.DB) map[string]handlerFunc {
	newListHandler := NewNewListHandler(bot)
	return map[string]handlerFunc{
		"newlist:0": newListHandler.ReplyCallback,
		"newlist:1": newListHandler.ReplyCallback1,
		// "newlist:~": newListHandler.ReplyCallback2,
	}
}

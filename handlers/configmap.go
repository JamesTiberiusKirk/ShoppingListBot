package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handlerFunc func(update tgbotapi.Update) error

func NewHandlerMap(bot *tgbotapi.BotAPI) map[string]handlerFunc {
	return map[string]handlerFunc{
		"newlist": NewNewListHandler(bot).Handle,
		"kb":      NewKeyboardHandler(bot).Handle,
	}
}

func GetHandlerCommandList() tgbotapi.SetMyCommandsConfig {
	return tgbotapi.NewSetMyCommands(
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

func NewCallbackMap(bot *tgbotapi.BotAPI) map[string]handlerFunc {
	return map[string]handlerFunc{
		"kb": NewKeyboardHandler(bot).Callback,
	}
}

package handlers

import (
	"errors"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/clients"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	JourneryExitErr = errors.New("exiting journery")
	UserErr         = errors.New("user error")
)

type JourneyTracker struct {
	Command     string
	Next        int
	PastUpdates []tgbotapi.Update
}

type HandlerFunc func(update tgbotapi.Update, previous []tgbotapi.Update) error
type HandlerInterface interface {
	// GetHandlerJourney returns handler funcs jouneys and weather or not the final elment in the array is to be called endlesly
	GetHandlerJourney() ([]HandlerFunc, bool)
}

func GetHandlerCommandList() tgbotapi.SetMyCommandsConfig {
	return tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Welcome and chat registration",
		},
		tgbotapi.BotCommand{
			Command:     "newlist",
			Description: "Create new shopping list",
		},
		// tgbotapi.BotCommand{
		// 	Command:     "kb",
		// 	Description: "Get keyboard",
		// },
	)
}

func NewHandlerJounreyMap(bot *tgbotapi.BotAPI, db *clients.DB) map[string]HandlerInterface {
	return map[string]HandlerInterface{
		"start":   NewStartHandler(bot.Send, db.AddNewChat, db.CheckIfChatExists),
		"newlist": NewNewListHandler(bot.Send, db.NewShoppingList),
	}
}

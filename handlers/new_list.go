package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewListHandler struct {
	bot *tgbotapi.BotAPI
}

func NewNewListHandler(bot *tgbotapi.BotAPI) *NewListHandler {
	return &NewListHandler{
		bot: bot,
	}
}

func (h *NewListHandler) Handle(update tgbotapi.Update) error {
	log.Print("New list called")
	return nil
}

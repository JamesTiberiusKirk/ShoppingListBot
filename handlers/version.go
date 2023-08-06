package handlers

import (
	"fmt"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/tgf"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type VersionHandler struct {
	sendMsg func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	getVer  func() string
}

func NewVersionHandler(sendMsg func(c tgbotapi.Chattable) (tgbotapi.Message, error), getVer func() string) *VersionHandler {
	return &VersionHandler{
		sendMsg: sendMsg,
		getVer:  getVer,
	}
}

func (h *VersionHandler) GetHandlerJourney() []tgf.HandlerFunc {
	return []tgf.HandlerFunc{
		func(ctx *tgf.Context) error {
			msg := tgbotapi.NewMessage(ctx.GetChatID(), fmt.Sprintf("ShoppingListsBot version: %s", h.getVer()))
			_, err := h.sendMsg(msg)
			if err != nil {
				ctx.Log.Error("Error sending message", "error", err)
				return err
			}
			return nil
		},
	}
}

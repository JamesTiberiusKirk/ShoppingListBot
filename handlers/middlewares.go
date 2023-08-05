package handlers

import (
	"fmt"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/tgf"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func chatRegistered(
	sendMsg func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	checkRegistration func(chatID int64) (bool, error),
	h tgf.HandlerFunc,
) tgf.HandlerFunc {
	return func(ctx *tgf.Context) error {
		ctx.Log.Info("[MIDDLEWARE]: checking chat registration")

		chatID := ctx.GetChatID()
		found, err := checkRegistration(chatID)
		if err != nil {
			ctx.Log.Error("Error could not check chat registration", "error", err)
			return fmt.Errorf("error looking for chat registration: %w", err)
		}
		if !found {
			msg := tgbotapi.NewMessage(chatID, "Chat not registered")
			_, err = sendMsg(msg)
			if err != nil {
				ctx.Log.Error("Error could not send message", "error", err)
				return fmt.Errorf("error senuding message: %w", err)
			}

			return tgf.UserErr
		}

		return h(ctx)
	}
}

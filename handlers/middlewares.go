package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func chatRegistered(
	sendMsg func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	checkRegistration func(chatID int64) (bool, error),
	h HandlerFunc,
) HandlerFunc {
	return func(context []byte, update tgbotapi.Update) (interface{}, error) {
		log.Print("[MIDDLEWARE]: Display List Handler")

		chatID, err := getChatID(update)
		if err != nil {
			return nil, fmt.Errorf("error finding chatID: %w", err)
		}

		found, err := checkRegistration(chatID)
		if err != nil {
			return nil, fmt.Errorf("error looking for chat registration: %w", err)
		}
		if !found {
			msg := tgbotapi.NewMessage(chatID, "Chat not registered")
			_, err = sendMsg(msg)
			if err != nil {
				return nil, fmt.Errorf("error senuding message: %w", err)
			}

			return nil, UserErr
		}

		return h(context, update)
	}
}

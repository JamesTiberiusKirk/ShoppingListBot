package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/inconshreveable/log15"
)

func chatRegistered(
	sendMsg func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	checkRegistration func(chatID int64) (bool, error),
	h HandlerFunc,
) HandlerFunc {
	return func(context []byte, update tgbotapi.Update) (interface{}, error) {
		log.Info("[MIDDLEWARE]: Display List Handler")

		chatID, err := getChatID(update)
		if err != nil {
			log.Error("Error could not get chat ID", "error", err)
			return nil, fmt.Errorf("error finding chatID: %w", err)
		}

		found, err := checkRegistration(chatID)
		if err != nil {
			log.Error("Error could not check chat registration", "error", err)
			return nil, fmt.Errorf("error looking for chat registration: %w", err)
		}
		if !found {
			msg := tgbotapi.NewMessage(chatID, "Chat not registered")
			_, err = sendMsg(msg)
			if err != nil {
				log.Error("Error could not send message", "error", err)
				return nil, fmt.Errorf("error senuding message: %w", err)
			}

			return nil, UserErr
		}

		return h(context, update)
	}
}

package handlers

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	messageNotFoundErr = errors.New("message not found")
)

func getMessage(update tgbotapi.Update) *tgbotapi.Message {
	if update.Message != nil {
		return update.Message
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Message
	}

	return nil
}

func getChatID(update tgbotapi.Update) (int64, error) {
	message := getMessage(update)
	if message == nil {
		return 0, messageNotFoundErr
	}

	return message.Chat.ID, nil
}

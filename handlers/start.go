package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
	sendMsg      func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	addChat      func(chatID int64) error
	checkIfExist func(chatID int64) (bool, error)
}

func NewStartHandler(sendMsg func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	addChat func(chatID int64) error,
	checkIfExist func(chatID int64) (bool, error)) *StartHandler {
	return &StartHandler{
		sendMsg:      sendMsg,
		addChat:      addChat,
		checkIfExist: checkIfExist,
	}
}

func (h *StartHandler) GetHandlerJourney() ([]HandlerFunc, bool) {
	return []HandlerFunc{
		func(context interface{}, update tgbotapi.Update) (interface{}, error) {

			log.Print("[HANDLER]: Start handler called")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the Telegram list manager, we are creating your account, bear with us.")
			_, err := h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			found, err := h.checkIfExist(update.Message.Chat.ID)
			if err != nil {
				log.Printf("[HANDLER ERROR]: when checking for existing chats: %s", err.Error())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred")
				_, err = h.sendMsg(msg)
				if err != nil {
					log.Printf("[HANDLER ERROR]: %s", err.Error())
					return nil, err
				}
				return nil, err
			}

			if found {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Chat already registered")
				_, err = h.sendMsg(msg)
				if err != nil {
					log.Printf("[HANDLER ERROR]: %s", err.Error())
					return nil, err
				}
				return nil, nil
			}

			err = h.addChat(update.Message.Chat.ID)
			if err != nil {
				log.Printf("[HANDLER ERROR]: %s", err.Error())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred")
				_, err = h.sendMsg(msg)
				if err != nil {
					log.Printf("[HANDLER ERROR]: %s", err.Error())
					return nil, err
				}
				return nil, err
			}

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Chat registered")
			_, err = h.sendMsg(msg)
			if err != nil {
				log.Printf("[HANDLER ERROR]: %s", err.Error())
				return nil, err
			}
			return nil, nil
		},
	}, false
}

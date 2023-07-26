package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/inconshreveable/log15"
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
		func(context []byte, update tgbotapi.Update) (interface{}, error) {

			log.Info("[HANDLER]: Start handler called")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the Telegram list manager, we are creating your account, bear with us.")
			_, err := h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			found, err := h.checkIfExist(update.Message.Chat.ID)
			if err != nil {
				log.Info("[HANDLER ERROR]: when checking for existing chats", "error", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred")
				_, err = h.sendMsg(msg)
				if err != nil {
					log.Error("Error sending message", "error", err)
					return nil, err
				}
				return nil, err
			}

			if found {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Chat already registered")
				_, err = h.sendMsg(msg)
				if err != nil {
					log.Error("Error sending message", "error", err)
					return nil, err
				}
				return nil, nil
			}

			err = h.addChat(update.Message.Chat.ID)
			if err != nil {
				log.Error("Error adding chat to db", "error", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred")
				_, err = h.sendMsg(msg)
				if err != nil {
					log.Error("Error sending message", "error", err)
					return nil, err
				}
				return nil, err
			}

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Chat registered")
			_, err = h.sendMsg(msg)
			if err != nil {
				log.Error("Error sending message", "error", err)
				return nil, err
			}
			return nil, nil
		},
	}, false
}

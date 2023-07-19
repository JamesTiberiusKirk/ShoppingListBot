package handlers

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewListHandler struct {
	sendMsg func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	addList func(chatID int64, title string, storeLoc string, dueDate *time.Time) error
}

func NewNewListHandler(
	msgSener func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	addList func(chatID int64, title string, storeLoc string, dueDate *time.Time) error,
) *NewListHandler {
	return &NewListHandler{
		sendMsg: msgSener,
		addList: addList,
	}
}

func (h *NewListHandler) GetHandlerJourney() ([]HandlerFunc, bool) {
	return []HandlerFunc{
		func(context interface{}, update tgbotapi.Update, previous []tgbotapi.Update) (interface{}, error) {
			log.Print("[HANDLER]: New List Handler")

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please Chose a name for the list")
			_, err := h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
		func(context interface{}, update tgbotapi.Update, previous []tgbotapi.Update) (interface{}, error) {
			log.Printf("[CALLBACK]: New list contextual reply callback with name %s", update.Message.Text)

			if update.Message.Text == "" {
				return nil, JourneryExitErr
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now, please chose a store")
			_, err := h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
		func(context interface{}, update tgbotapi.Update, previous []tgbotapi.Update) (interface{}, error) {
			log.Printf("[CALLBACK]: New list contextual reply callback 1 with name %s", update.Message.Text)

			if update.Message.Text == "" {
				return nil, JourneryExitErr
			}

			chatID := update.Message.Chat.ID
			name := previous[1].Message.Text
			store := update.Message.Text

			err := h.addList(chatID, name, store, nil)
			if err != nil {
				return nil, fmt.Errorf("error inserting shopping_list: %w", err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "List created, thank you")
			_, err = h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}, false
}

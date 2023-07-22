package handlers

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewListHandler struct {
	sendMsg           func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	addList           func(chatID int64, title string, storeLoc string, dueDate *time.Time) error
	checkRegistration func(chatID int64) (bool, error)
}

func NewNewListHandler(
	msgSener func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	addList func(chatID int64, title string, storeLoc string, dueDate *time.Time) error,
	checkRegistration func(chatID int64) (bool, error),
) *NewListHandler {
	return &NewListHandler{
		sendMsg:           msgSener,
		addList:           addList,
		checkRegistration: checkRegistration,
	}
}

type NewListHandlerContext struct {
	Title   string
	Store   string
	DueDate *time.Time
}

// TODO: Refactor this to not use previous then remove previousfrom the entire application
func (h *NewListHandler) GetHandlerJourney() ([]HandlerFunc, bool) {
	return []HandlerFunc{
		chatRegistered(h.sendMsg, h.checkRegistration,
			func(context interface{}, update tgbotapi.Update) (interface{}, error) {
				log.Print("[HANDLER]: New List Handler")

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please Chose a name for the list")
				_, err := h.sendMsg(msg)
				if err != nil {
					return nil, err
				}

				c := NewListHandlerContext{}

				return c, nil
			},
		),
		func(context interface{}, update tgbotapi.Update) (interface{}, error) {
			log.Printf("[CALLBACK]: New list contextual reply callback with name %s", update.Message.Text)

			chatID, _ := getChatID(update)

			if update.Message.Text == "" {
				return nil, JourneryExitErr
			}

			c, ok := context.(NewListHandlerContext)
			if !ok {
				return nil, CouldNotExteactContextErr
			}

			c.Title = update.Message.Text
			msg := tgbotapi.NewMessage(chatID, "Now, please chose a store")
			_, err := h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		func(context interface{}, update tgbotapi.Update) (interface{}, error) {
			log.Printf("[CALLBACK]: New list contextual reply callback 1 with name %s", update.Message.Text)

			if update.Message.Text == "" {
				return nil, JourneryExitErr
			}

			c, ok := context.(NewListHandlerContext)
			if !ok {
				return nil, CouldNotExteactContextErr
			}
			chatID, _ := getChatID(update)
			c.Store = update.Message.Text

			err := h.addList(chatID, c.Title, c.Store, c.DueDate)
			if err != nil {
				return nil, fmt.Errorf("error inserting shopping_list: %w", err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "List created, thank you")
			_, err = h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}, false
}

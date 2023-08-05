package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/tgf"
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
func (h *NewListHandler) GetHandlerJourney() []tgf.HandlerFunc {
	return []tgf.HandlerFunc{
		chatRegistered(h.sendMsg, h.checkRegistration,
			func(ctx *tgf.Context) error {
				ctx.Log.Info("[HANDLER]: New List Handler")

				msg := tgbotapi.NewMessage(ctx.GetChatID(), "Please Chose a name for the list")
				_, err := h.sendMsg(msg)
				if err != nil {
					ctx.Log.Error("Error sending message %w", err)
					return err
				}

				c := NewListHandlerContext{}
				return ctx.SetContexData(c)
			},
		),
		func(ctx *tgf.Context) error {
			message := ctx.GetMessage()
			chatID := ctx.GetChatID()
			ctx.Log.Info("[HANDLER]: New list contextual reply callback with name %s", message.Text)

			if message.Text == "" {
				ctx.Exit()
				return nil
			}

			var c NewListHandlerContext
			err := json.Unmarshal(ctx.Journey.RawContext, &c)
			if err != nil {
				ctx.Log.Error("Error unmarshaling context", "error", err)
				return fmt.Errorf("%w: %w", tgf.CouldNotExteactContextErr, err)
			}

			c.Title = message.Text
			msg := tgbotapi.NewMessage(chatID, "Now, please chose a store")
			_, err = h.sendMsg(msg)
			if err != nil {
				ctx.Log.Error("Error sending message", "error", err)
				return err
			}

			return ctx.SetContexData(c)
		},
		func(ctx *tgf.Context) error {
			message := ctx.GetMessage()

			ctx.Log.Info("[HANDLER]: New list contextual reply callback 1 with name", message.Text)

			if message.Text == "" {
				ctx.Exit()
				return nil
			}

			var c NewListHandlerContext
			err := json.Unmarshal(ctx.Journey.RawContext, &c)
			if err != nil {
				ctx.Log.Error("Error unmarshaling context", "error", err)
				return fmt.Errorf("%w: %w", tgf.CouldNotExteactContextErr, err)
			}

			chatID := ctx.GetChatID()
			c.Store = message.Text

			err = h.addList(chatID, c.Title, c.Store, c.DueDate)
			if err != nil {
				ctx.Log.Error("Error inserting shopping list", "error", err)
				return fmt.Errorf("error inserting shopping_list: %w", err)
			}

			msg := tgbotapi.NewMessage(message.Chat.ID, "List created, thank you")
			_, err = h.sendMsg(msg)
			if err != nil {
				return err
			}

			return ctx.SetContexData(c)
		},
	}
}

package handlers

import (
	"github.com/JamesTiberiusKirk/tgf"
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

func (h *StartHandler) GetHandlerJourney() []tgf.HandlerFunc {
	return []tgf.HandlerFunc{
		func(ctx *tgf.Context) error {
			ctx.Log.Info("[HANDLER]: Start handler called")
			chatID := ctx.GetChatID()

			msg := tgbotapi.NewMessage(chatID, "Welcome to the Telegram list manager, we are creating your account, bear with us.")
			_, err := h.sendMsg(msg)
			if err != nil {
				return err
			}

			found, err := h.checkIfExist(chatID)
			if err != nil {
				ctx.Log.Info("[HANDLER ERROR]: when checking for existing chats", "error", err)
				msg := tgbotapi.NewMessage(chatID, "An error occurred")
				_, err = h.sendMsg(msg)
				if err != nil {
					ctx.Log.Error("Error sending message", "error", err)
					return err
				}
				return err
			}

			if found {
				msg = tgbotapi.NewMessage(chatID, "Chat already registered")
				_, err = h.sendMsg(msg)
				if err != nil {
					ctx.Log.Error("Error sending message", "error", err)
					return err
				}
				return nil
			}

			err = h.addChat(chatID)
			if err != nil {
				ctx.Log.Error("Error adding chat to db", "error", err)
				msg := tgbotapi.NewMessage(chatID, "An error occurred")
				_, err = h.sendMsg(msg)
				if err != nil {
					ctx.Log.Error("Error sending message", "error", err)
					return err
				}
				return err
			}

			msg = tgbotapi.NewMessage(chatID, "Chat registered")
			_, err = h.sendMsg(msg)
			if err != nil {
				ctx.Log.Error("Error sending message", "error", err)
				return err
			}

			return nil
		},
	}
}

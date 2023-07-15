package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
	bot          *tgbotapi.BotAPI
	addChat      func(chatID int64) error
	checkIfExist func(chatID int64) (bool, error)
}

func NewStartHandler(bot *tgbotapi.BotAPI, addChat func(chatID int64) error, checkIfExist func(chatID int64) (bool, error)) *StartHandler {
	return &StartHandler{
		bot:          bot,
		addChat:      addChat,
		checkIfExist: checkIfExist,
	}
}

func (h *StartHandler) Handle(update tgbotapi.Update) error {
	log.Print("[HANDLER]: Start handler called")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the Telegram list manager, we are creating your account, bear with us.")
	_, err := h.bot.Send(msg)
	if err != nil {
		return err
	}

	found, err := h.checkIfExist(update.Message.Chat.ID)
	if err != nil {
		log.Printf("[HANDLER ERROR]: when checking for existing chats: %s", err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Printf("[HANDLER ERROR]: %s", err.Error())
			return err
		}
		return err
	}

	if found {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Chat already registered")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Printf("[HANDLER ERROR]: %s", err.Error())
			return err
		}
		return nil
	}

	err = h.addChat(update.Message.Chat.ID)
	if err != nil {
		log.Printf("[HANDLER ERROR]: %s", err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "An error occurred")
		_, err = h.bot.Send(msg)
		if err != nil {
			log.Printf("[HANDLER ERROR]: %s", err.Error())
			return err
		}
		return err
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Chat registered")
	_, err = h.bot.Send(msg)
	if err != nil {
		log.Printf("[HANDLER ERROR]: %s", err.Error())
		return err
	}

	return nil
}

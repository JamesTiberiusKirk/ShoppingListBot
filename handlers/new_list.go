package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NewListHandler struct {
	bot     *tgbotapi.BotAPI
	AddList func(chatID int64, store string, name string) error
}

func NewNewListHandler(bot *tgbotapi.BotAPI) *NewListHandler {
	return &NewListHandler{
		bot: bot,
	}
}

func (h *NewListHandler) Handle(update tgbotapi.Update) error {
	log.Print("[HANDLER]: New List Handler")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	msg.Text = "Please Chose a name for the list"

	// Send the message.
	_, err := h.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (h *NewListHandler) ReplyCallback(update tgbotapi.Update) error {
	log.Printf("[CALLBACK]: New list contextual reply callback with name %s", update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	msg.Text = "Now, please chose a store"

	// Send the message.
	_, err := h.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (h *NewListHandler) ReplyCallback1(update tgbotapi.Update) error {
	log.Printf("[CALLBACK]: New list contextual reply callback 1 with name %s", update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	msg.Text = "Thank you"

	// Send the message.
	_, err := h.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (h *NewListHandler) ReplyCallback2(update tgbotapi.Update) error {
	log.Printf("[CALLBACK]: New list contextual reply callback 2 with name %s", update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	msg.Text = "Thank you again"

	// Send the message.
	_, err := h.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

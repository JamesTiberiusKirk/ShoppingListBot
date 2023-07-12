package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var kbData = map[string]string{
	"apples":        "Apples",
	"oranges":       "Oranges",
	"beef-mince-%5": "Beef mince %5",
}

func buildKeyboardInline() tgbotapi.InlineKeyboardMarkup {

	rows := [][]tgbotapi.InlineKeyboardButton{}
	for k, v := range kbData {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(v, k)))
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(
		rows...,
	)
	return kb
}

func buildKeyboard() tgbotapi.ReplyKeyboardMarkup {
	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("This is a test1"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("This is a test2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("This is a test3"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("This is a test4"),
		),
	)
	return kb
}

// var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
// 		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
// 		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
// 	),
// 	tgbotapi.NewInlineKeyboardRow(
// 		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
// 		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
// 		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
// 	),
// )

type KeyboardHandler struct {
	bot *tgbotapi.BotAPI
}

func NewKeyboardHandler(bot *tgbotapi.BotAPI) *KeyboardHandler {
	return &KeyboardHandler{
		bot: bot,
	}
}

func (h *KeyboardHandler) Handle(update tgbotapi.Update) error {
	// Construct a new message from the given chat ID and containing
	// the text that we received.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	// If the message was open, add a copy of our numeric keyboard.
	// msg.ReplyMarkup = numericKeyboard
	msg.ReplyMarkup = buildKeyboardInline()

	// Send the message.
	_, err := h.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (h *KeyboardHandler) Callback(update tgbotapi.Update) error {
	log.Printf("[CALLBACK]: %s, DATA: %s", update.CallbackQuery.ID, update.CallbackQuery.Data)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := h.bot.Request(callback); err != nil {
		return err
	}

	kbData[update.CallbackQuery.Data] = fmt.Sprintf("✅ %s", kbData[update.CallbackQuery.Data])
	req := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, buildKeyboardInline())
	h.bot.Request(req)

	return nil
}

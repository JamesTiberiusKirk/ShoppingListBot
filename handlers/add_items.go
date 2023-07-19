package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddItemsHandler struct {
	ShoppingListsMap map[string]types.ShoppingList
	ShoppingList     types.ShoppingList
	Items            []string
	sendMsg          func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	getLists         func(chatID int64) ([]types.ShoppingList, error)
	addItems         func(listID string, itemText []string) error
}

func NewAddItemsHandler(
	msgSener func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	getLists func(chatID int64) ([]types.ShoppingList, error),
	addItems func(listID string, itemText []string) error,
) *AddItemsHandler {
	return &AddItemsHandler{
		sendMsg:  msgSener,
		getLists: getLists,
		addItems: addItems,
	}
}

type AddItemsHandlerContext struct {
	ShoppingListsMap map[string]types.ShoppingList
	ShoppingList     types.ShoppingList
	Items            []string
}

func (h *AddItemsHandler) GetHandlerJourney() ([]HandlerFunc, bool) {
	return []HandlerFunc{
		func(update tgbotapi.Update, previous []tgbotapi.Update) error {
			log.Print("[HANDLER]: Add Items Handler")

			lists, err := h.getLists(update.Message.Chat.ID)
			if err != nil {
				return err
			}

			kbRows := [][]tgbotapi.InlineKeyboardButton{}
			h.ShoppingListsMap = map[string]types.ShoppingList{}
			for _, l := range lists {
				kbRows = append(
					kbRows,
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s - %s", l.Title, l.StoreLocation), l.ID),
					),
				)
				h.ShoppingListsMap[l.ID] = l
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please chose the list to add items to")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(kbRows...)
			_, err = h.sendMsg(msg)
			if err != nil {
				return err
			}

			return nil
		},
		func(update tgbotapi.Update, previous []tgbotapi.Update) error {
			log.Print("[HANDLER]: Add Items Handler 2")
			chatID := update.CallbackQuery.Message.Chat.ID
			listID := update.CallbackQuery.Data

			msg := tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf("Adding to %s, start typing the items and type \"DONE\" when finished", listID),
			)

			_, err := h.sendMsg(msg)
			if err != nil {
				return err
			}

			h.ShoppingList = h.ShoppingListsMap[listID]

			return nil
		},
		func(update tgbotapi.Update, previous []tgbotapi.Update) error {
			log.Print("[HANDLER]: Add Items Handler 2")

			var message tgbotapi.Message
			if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
				message = *update.CallbackQuery.Message
			} else if update.Message != nil {
				message = *update.Message
			} else {
				return JourneryExitErr
			}

			if strings.ToUpper(message.Text) == "DONE" {
				log.Printf("[HANDLER]: ITEMS: %+v", h.Items)
				err := h.addItems(h.ShoppingList.ID, h.Items)
				if err != nil {
					return fmt.Errorf("error inserting items into list %s, %w", h.ShoppingList.ID, err)
				}

				return JourneryExitErr
			}

			h.Items = append(h.Items, message.Text)

			return nil
		},
	}, true
}

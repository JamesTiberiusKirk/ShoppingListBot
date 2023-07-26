package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddItemsHandler struct {
	sendMsg  func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	getLists func(chatID int64) ([]types.ShoppingList, error)
	addItems func(listID string, itemText []string) error
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
		func(context []byte, update tgbotapi.Update) (interface{}, error) {
			log.Print("[HANDLER]: Add Items Handler")

			chatID, _ := getChatID(update)
			lists, err := h.getLists(update.Message.Chat.ID)
			if err != nil {
				return nil, err
			}

			if len(lists) < 1 {
				msg := tgbotapi.NewMessage(chatID, "There are no lists")
				_, err = h.sendMsg(msg)
				if err != nil {
					return nil, err
				}
				return nil, JourneryExitErr
			}

			c := AddItemsHandlerContext{}

			kbRows := [][]tgbotapi.InlineKeyboardButton{}
			c.ShoppingListsMap = map[string]types.ShoppingList{}
			for _, l := range lists {
				kbRows = append(
					kbRows,
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s - %s", l.Title, l.StoreLocation), l.ID),
					),
				)
				c.ShoppingListsMap[l.ID] = l
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please chose the list to add items to")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(kbRows...)
			_, err = h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		func(context []byte, update tgbotapi.Update) (interface{}, error) {
			log.Print("[HANDLER]: Add Items Handler 2")
			chatID := update.CallbackQuery.Message.Chat.ID
			listID := update.CallbackQuery.Data

			var c AddItemsHandlerContext
			err := json.Unmarshal(context, &c)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", CouldNotExteactContextErr, err)
			}

			c.ShoppingList = c.ShoppingListsMap[listID]

			msg := tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf("Adding to %s, start typing the items and type \"DONE\" when finished", c.ShoppingList.Title),
			)

			_, err = h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		func(context []byte, update tgbotapi.Update) (interface{}, error) {
			log.Print("[HANDLER]: Add Items Handler 2")
			chatID, _ := getChatID(update)

			var message tgbotapi.Message
			if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
				message = *update.CallbackQuery.Message
			} else if update.Message != nil {
				message = *update.Message
			} else {
				return nil, JourneryExitErr
			}

			var c AddItemsHandlerContext
			err := json.Unmarshal(context, &c)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", CouldNotExteactContextErr, err)
			}

			if strings.ToUpper(message.Text) == "DONE" {
				log.Printf("[HANDLER]: ITEMS: %+v", c.Items)
				err := h.addItems(c.ShoppingList.ID, c.Items)
				if err != nil {
					return nil, fmt.Errorf("error inserting items into list %s, %w", c.ShoppingList.ID, err)
				}

				// TODO: need to add some nice formatting here so im not just dumping the array on the suer
				msg := tgbotapi.NewMessage(
					chatID,
					fmt.Sprintf("Adding items %+v, to %s", c.Items, c.ShoppingList.Title),
				)

				_, err = h.sendMsg(msg)
				if err != nil {
					return nil, err
				}

				return nil, JourneryExitErr
			}

			c.Items = append(c.Items, message.Text)

			return c, nil
		},
	}, true
}

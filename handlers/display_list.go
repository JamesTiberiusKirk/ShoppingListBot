package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DisplayListHandler struct {
	sendMsg            func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	botReq             func(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	getLists           func(chatID int64) ([]types.ShoppingList, error)
	getItems           func(listID string) ([]types.ShoppingListItem, error)
	toggleItemPurchase func(itemID string) error
	checkRegistration  func(chatID int64) (bool, error)
}

func NewDisplayListHandler(
	msgSener func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	botReq func(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error),
	getLists func(chatID int64) ([]types.ShoppingList, error),
	getItems func(listID string) ([]types.ShoppingListItem, error),
	toggleItemPurchase func(itemID string) error,
	checkRegistration func(chatID int64) (bool, error),
) *DisplayListHandler {
	return &DisplayListHandler{
		sendMsg:            msgSener,
		botReq:             botReq,
		getLists:           getLists,
		getItems:           getItems,
		toggleItemPurchase: toggleItemPurchase,
		checkRegistration:  checkRegistration,
	}
}

type DisplayListHandlerContext struct {
	ShoppingListsMap map[string]types.ShoppingList
	ShoppingList     types.ShoppingList
	Items            []types.ShoppingListItem
}

func (h *DisplayListHandler) GetHandlerJourney() ([]HandlerFunc, bool) {
	return []HandlerFunc{
		chatRegistered(h.sendMsg, h.checkRegistration,
			func(context []byte, update tgbotapi.Update) (interface{}, error) {
				log.Print("[HANDLER]: Display List Handler")

				lists, err := h.getLists(update.Message.Chat.ID)
				if err != nil {
					return nil, err
				}

				if len(lists) < 1 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There are no lists to chose from")
					_, err = h.sendMsg(msg)
					if err != nil {
						return nil, err
					}
					return nil, JourneryExitErr
				}

				c := DisplayListHandlerContext{
					ShoppingListsMap: map[string]types.ShoppingList{},
					Items:            []types.ShoppingListItem{},
				}

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

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please chose the list to display")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(kbRows...)
				_, err = h.sendMsg(msg)
				if err != nil {
					return nil, err
				}

				return c, nil
			}),
		func(context []byte, update tgbotapi.Update) (interface{}, error) {
			log.Print("[HANDLER]: Display List Handler 2")

			var c DisplayListHandlerContext
			err := json.Unmarshal(context, &c)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", CouldNotExteactContextErr, err)
			}

			listID := update.CallbackQuery.Data
			c.ShoppingList = c.ShoppingListsMap[listID]

			items, err := h.getItems(listID)
			if err != nil {
				return nil, fmt.Errorf("error getting items from db: %w", err)
			}

			if len(items) < 1 {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("There are no items in list %s", c.ShoppingList.Title))
				_, err = h.sendMsg(msg)
				if err != nil {
					return nil, err
				}
				return nil, JourneryExitErr
			}

			for _, i := range items {
				c.Items = append(c.Items, i)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please chose the list to display")
			msg.ReplyMarkup = buildItemsKeyboard(c)
			_, err = h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		func(context []byte, update tgbotapi.Update) (interface{}, error) {
			log.Print("[HANDLER]: Display List Handler 2")

			var c DisplayListHandlerContext
			err := json.Unmarshal(context, &c)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", CouldNotExteactContextErr, err)
			}

			itemID := ""
			data := update.CallbackQuery.Data
			switch data {
			case "edit":
			case "done":

			default:
				itemID = data
			}

			itemIndex := -1
			for index, i := range c.Items {
				if i.ID == itemID {
					itemIndex = index
				}
			}

			if itemIndex == -1 {
				return nil, fmt.Errorf("could not find item ID: %s", itemID)
			}

			err = h.toggleItemPurchase(c.Items[itemIndex].ID)
			if err != nil {
				return nil, fmt.Errorf("error toggling item purchace in db id: %s, err: %w", c.Items[itemIndex].ID, err)
			}
			// NOTE: on a technical level this could present a race condition since it does not display db values
			// but since this can only be modified on one chat then should be fine
			c.Items[itemIndex].Purchased = !c.Items[itemIndex].Purchased

			markup := buildItemsKeyboard(c)
			msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, markup)
			_, err = h.botReq(msg)
			if err != nil {
				return nil, fmt.Errorf("error making bot request: %w", err)
			}

			return c, nil
		},
	}, true
}

func buildItemsKeyboard(c DisplayListHandlerContext) tgbotapi.InlineKeyboardMarkup {
	kbRows := [][]tgbotapi.InlineKeyboardButton{}
	for _, i := range c.Items {
		text := ""
		if i.Purchased {
			log.Print("PURCHASED")
			text += "✅ "
		}
		text += i.ItemText
		log.Print(text)

		kbRows = append(
			kbRows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(text, i.ID),
			),
		)
		// TODO: need to make another bottom KB row in order to allow the user to exit or modify the list
	}

	kbRows = append(
		kbRows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
			tgbotapi.NewInlineKeyboardButtonData("Edit", "edit"),
		),
	)

	return tgbotapi.NewInlineKeyboardMarkup(kbRows...)
}

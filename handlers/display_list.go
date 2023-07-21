package handlers

import (
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
			func(context interface{}, update tgbotapi.Update, previous []tgbotapi.Update) (interface{}, error) {
				log.Print("[HANDLER]: Display List Handler")

				lists, err := h.getLists(update.Message.Chat.ID)
				if err != nil {
					return nil, err
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
		func(context interface{}, update tgbotapi.Update, previous []tgbotapi.Update) (interface{}, error) {
			log.Print("[HANDLER]: Display List Handler 2")
			c, ok := context.(DisplayListHandlerContext)
			if !ok {
				return nil, CouldNotExteactContextErr
			}

			listID := update.CallbackQuery.Data
			c.ShoppingList = c.ShoppingListsMap[listID]

			items, err := h.getItems(listID)
			if err != nil {
				return nil, fmt.Errorf("error getting items from db: %w", err)
			}

			for _, i := range items {
				c.Items = append(c.Items, i)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please chose the list to display")
			msg.ReplyMarkup = h.buildKeyboard(c)
			_, err = h.sendMsg(msg)
			if err != nil {
				return nil, err
			}

			return c, nil
		},
		func(context interface{}, update tgbotapi.Update, previous []tgbotapi.Update) (interface{}, error) {
			log.Print("[HANDLER]: Display List Handler 2")

			c, ok := context.(DisplayListHandlerContext)
			if !ok {
				return nil, CouldNotExteactContextErr
			}

			itemID := update.CallbackQuery.Data
			itemIndex := -1
			for index, i := range c.Items {
				if i.ID == itemID {
					itemIndex = index
				}
			}

			if itemIndex == -1 {
				return nil, fmt.Errorf("could not find item ID: %s", itemID)
			}

			log.Print(c.Items[itemIndex].Purchased)
			c.Items[itemIndex].Purchased = !c.Items[itemIndex].Purchased
			log.Print(c.Items[itemIndex].Purchased)

			err := h.toggleItemPurchase(c.Items[itemIndex].ID)
			if err != nil {
				return nil, fmt.Errorf("error toggling item purchace in db id: %s, err: %w", c.Items[itemIndex].ID, err)
			}

			markup := h.buildKeyboard(c)
			msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, markup)
			_, err = h.botReq(msg)
			if err != nil {
				return nil, fmt.Errorf("error making bot request: %w", err)
			}

			return c, nil
		},
	}, true
}

func (h *DisplayListHandler) buildKeyboard(c DisplayListHandlerContext) tgbotapi.InlineKeyboardMarkup {
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
	return tgbotapi.NewInlineKeyboardMarkup(kbRows...)
}

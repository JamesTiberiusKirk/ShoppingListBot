package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/inconshreveable/log15"
)

type DisplayListHandler struct {
	sendMsg            func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	botReq             func(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	getLists           func(chatID int64) ([]types.ShoppingList, error)
	getItems           func(listID string, togglePurchased bool) ([]types.ShoppingListItem, error)
	toggleItemPurchase func(itemID string) error
	checkRegistration  func(chatID int64) (bool, error)
	deleteItem         func(itemID string) error
}

func NewDisplayListHandler(
	msgSener func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	botReq func(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error),
	getLists func(chatID int64) ([]types.ShoppingList, error),
	getItems func(listID string, togglePurchased bool) ([]types.ShoppingListItem, error),
	toggleItemPurchase func(itemID string) error,
	checkRegistration func(chatID int64) (bool, error),
	deleteItem func(itemID string) error,
) *DisplayListHandler {
	return &DisplayListHandler{
		sendMsg:            msgSener,
		botReq:             botReq,
		getLists:           getLists,
		getItems:           getItems,
		toggleItemPurchase: toggleItemPurchase,
		checkRegistration:  checkRegistration,
		deleteItem:         deleteItem,
	}
}

type DisplayListHandlerContext struct {
	ShoppingListsMap      map[string]types.ShoppingList
	ShoppingList          types.ShoppingList
	Items                 []types.ShoppingListItem
	ItemsKeyboardEditable bool
	ShowPurchasedItems    bool
}

func (h *DisplayListHandler) GetHandlerJourney() ([]HandlerFunc, bool) {
	return []HandlerFunc{
		chatRegistered(h.sendMsg, h.checkRegistration,
			func(context []byte, update tgbotapi.Update) (interface{}, error) {
				log.Info("[HANDLER]: Display List Handler")

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
					ShoppingListsMap:   map[string]types.ShoppingList{},
					Items:              []types.ShoppingListItem{},
					ShowPurchasedItems: false,
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
					log.Error("Error sending message", "error", err)
					return nil, err
				}

				return c, nil
			}),
		func(context []byte, update tgbotapi.Update) (interface{}, error) {
			log.Info("[HANDLER]: Display List Handler 2")

			var c DisplayListHandlerContext
			err := json.Unmarshal(context, &c)
			if err != nil {
				log.Error("Could not extract context", "error", err)
				return nil, fmt.Errorf("%w: %w", CouldNotExteactContextErr, err)
			}

			listID := update.CallbackQuery.Data
			c.ShoppingList = c.ShoppingListsMap[listID]

			items, err := h.getItems(listID, c.ShowPurchasedItems)
			if err != nil {
				log.Error("Error getting items", "error", err)
				return nil, fmt.Errorf("error getting items from db: %w", err)
			}

			if len(items) < 1 {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("There are no items in list %s", c.ShoppingList.Title))
				_, err = h.sendMsg(msg)
				if err != nil {
					log.Error("Error sending message", "error", err)
					return nil, err
				}
				return nil, JourneryExitErr
			}

			for _, i := range items {
				c.Items = append(c.Items, i)
			}

			markup := h.buildItemsKeyboard(c)
			msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Please chose the list to display", markup)
			_, err = h.botReq(msg)
			if err != nil {
				log.Error("Error sending bot request", "error", err)
				return nil, fmt.Errorf("error making bot request: %w", err)
			}

			return c, nil
		},
		func(context []byte, update tgbotapi.Update) (interface{}, error) {
			log.Info("[HANDLER]: Display List Handler 3")

			var c DisplayListHandlerContext
			err := json.Unmarshal(context, &c)
			if err != nil {
				log.Error("Error unmarshaling context", "error", err)
				return nil, fmt.Errorf("%w: %w", CouldNotExteactContextErr, err)
			}

			itemID := ""
			data := update.CallbackQuery.Data

			splitData := strings.Split(data, ":")

			// TODO: Posibly refactor, not the happiest with this but ok for now
			switch splitData[0] {
			case "del":
				toDeleteID := splitData[1]
				log.Info("Deleting", "item", toDeleteID)
				err := h.deleteItem(toDeleteID)
				if err != nil {
					return nil, fmt.Errorf("error deleting item from db: %s, err: %w", toDeleteID, err)
				}

				itemIndex := -1
				for index, i := range c.Items {
					if i.ID == toDeleteID {
						itemIndex = index
					}
				}

				if itemIndex == -1 {
					log.Error("Error could not find itemIndex")
					return nil, fmt.Errorf("could not find item ID: %s", itemID)
				}

				c.Items = append(c.Items[:itemIndex], c.Items[itemIndex+1:]...)
			case "togglePurchased":
				c.ShowPurchasedItems = !c.ShowPurchasedItems
				fmt.Printf("%+v", c)
				items, err := h.getItems(c.ShoppingList.ID, c.ShowPurchasedItems)
				if err != nil {
					log.Error("Error getting items", "error", err)
					return nil, fmt.Errorf("error getting items from db: %w", err)
				}

				c.Items = items
			case "edit":
				c.ItemsKeyboardEditable = !c.ItemsKeyboardEditable
				log.Info("setting editable", "editable", c.ItemsKeyboardEditable)
			case "done":
				deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				_, err = h.botReq(deleteMsg)
				if err != nil {
					log.Error("Error deleting inline keyboard", "error", err)
					return nil, fmt.Errorf("error making bot request: %w", err)
				}
				return nil, JourneryExitErr
			default:
				itemID = data

				itemIndex := -1
				for index, i := range c.Items {
					if i.ID == itemID {
						itemIndex = index
					}
				}

				if itemIndex == -1 {
					log.Error("Error could not find itemIndex")
					return nil, fmt.Errorf("could not find item ID: %s", itemID)
				}

				err = h.toggleItemPurchase(c.Items[itemIndex].ID)
				if err != nil {
					log.Error("Error toggling item purchace", "error", err)
					return nil, fmt.Errorf("error toggling item purchace in db id: %s, err: %w", c.Items[itemIndex].ID, err)
				}
				c.Items[itemIndex].Purchased = !c.Items[itemIndex].Purchased
			}

			markup := h.buildItemsKeyboard(c)
			msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, markup)
			_, err = h.botReq(msg)
			if err != nil {
				log.Error("Error sending bot request", "error", err)
				return nil, fmt.Errorf("error making bot request: %w", err)
			}

			return c, nil
		},
	}, true
}

func (h *DisplayListHandler) buildItemsKeyboard(c DisplayListHandlerContext) tgbotapi.InlineKeyboardMarkup {
	kbRows := [][]tgbotapi.InlineKeyboardButton{}
	for _, i := range c.Items {
		text := ""
		if i.Purchased {
			text += "✅ "
		}
		text += i.ItemText

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, i.ID),
		)

		if c.ItemsKeyboardEditable {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData("❌", "del:"+i.ID))
		}

		kbRows = append(kbRows, row)
	}

	kbRows = append(
		kbRows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Show Purchased", "togglePurchased"),
			tgbotapi.NewInlineKeyboardButtonData("Edit", "edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
		),
	)

	return tgbotapi.NewInlineKeyboardMarkup(kbRows...)
}

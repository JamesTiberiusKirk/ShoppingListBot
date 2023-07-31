package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/inconshreveable/log15"
)

type EditListsHandler struct {
	sendMsg           func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	botReq            func(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	getLists          func(chatID int64) ([]types.ShoppingList, error)
	getItems          func(listID string, togglePurchased bool) ([]types.ShoppingListItem, error)
	checkRegistration func(chatID int64) (bool, error)
	deleteList        func(listID string) error
}

func NewEditListsHandler(
	msgSener func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	botReq func(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error),
	getLists func(chatID int64) ([]types.ShoppingList, error),
	getItems func(listID string, togglePurchased bool) ([]types.ShoppingListItem, error),
	checkRegistration func(chatID int64) (bool, error),
	deleteList func(itemID string) error,
) *EditListsHandler {
	return &EditListsHandler{
		sendMsg:  msgSener,
		botReq:   botReq,
		getLists: getLists,
		// getItems:          getItems,
		checkRegistration: checkRegistration,
		deleteList:        deleteList,
	}
}

type EditListsHandlerContext struct {
	SelectedList  string
	ShoppingLists []types.ShoppingList
	Items         []types.ShoppingListItem
}

func (h *EditListsHandler) GetHandlerJourney() ([]HandlerFunc, bool) {
	return []HandlerFunc{
		chatRegistered(h.sendMsg, h.checkRegistration,
			func(context []byte, update tgbotapi.Update) (interface{}, error) {
				log.Info("[HANDLER]: Display List Handler")

				c := EditListsHandlerContext{
					ShoppingLists: []types.ShoppingList{},
					Items:         []types.ShoppingListItem{},
				}

				if context != nil {
					err := json.Unmarshal(context, &c)
					if err != nil {
						log.Error("Error unmarshaling context", "error", err)
						return nil, fmt.Errorf("%w: %w", CouldNotExteactContextErr, err)
					}
				}

				chatID, _ := getChatID(update)

				data := ""
				if update.CallbackQuery != nil {
					data = update.CallbackQuery.Data
				}

				switch data {
				case "delete":
					log.Info("deleting", "selected_list", c.SelectedList)

					if c.SelectedList == "" {
						return nil, nil
					}

					err := h.deleteList(c.SelectedList)
					if err != nil {
						log.Error("Error deleting list", "error", err)
						return nil, fmt.Errorf("error deleting list from db: %w", err)
					}
					lists, err := h.getLists(chatID)
					if err != nil {
						return nil, err
					}
					c.ShoppingLists = lists

				case "done":
					deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
					_, err := h.botReq(deleteMsg)
					if err != nil {
						log.Error("Error deleting inline keyboard", "error", err)
						return nil, fmt.Errorf("error making bot request: %w", err)
					}
					return nil, JourneryExitErr
				default:
					if data != "" {
						c.SelectedList = data
					}

					lists, err := h.getLists(chatID)
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
					c.ShoppingLists = lists
				}

				kbRows := h.buildKeyboard(c)

				if update.Message != nil {
					msg := tgbotapi.NewMessage(chatID, "Please chose the list to edit")
					msg.ReplyMarkup = kbRows
					_, err := h.sendMsg(msg)
					if err != nil {
						log.Error("Error sending message", "error", err)
						return nil, err
					}
				} else {
					msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, kbRows)
					_, err := h.botReq(msg)
					if err != nil {
						log.Error("Error sending bot request", "error", err)
						return nil, fmt.Errorf("error making bot request: %w", err)
					}
				}

				return c, nil
			}),
	}, true
}

func (h *EditListsHandler) buildKeyboard(c EditListsHandlerContext) tgbotapi.InlineKeyboardMarkup {
	kbRows := [][]tgbotapi.InlineKeyboardButton{}
	for _, l := range c.ShoppingLists {

		var row []tgbotapi.InlineKeyboardButton
		if l.ID == c.SelectedList {
			row = tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("👉 %s - %s", l.Title, l.StoreLocation), l.ID),
			)
		} else {
			row = tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s - %s", l.Title, l.StoreLocation), l.ID),
			)
		}

		kbRows = append(
			kbRows,
			row,
		)
	}

	kbRows = append(
		kbRows,
		tgbotapi.NewInlineKeyboardRow(
			// tgbotapi.NewInlineKeyboardButtonData("Edit", "edit"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Delete", "delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
		),
	)

	return tgbotapi.NewInlineKeyboardMarkup(kbRows...)
}

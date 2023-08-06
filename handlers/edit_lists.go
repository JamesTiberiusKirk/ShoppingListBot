package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/tgf"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func (h *EditListsHandler) GetHandlerJourney() []tgf.HandlerFunc {
	return []tgf.HandlerFunc{
		chatRegistered(h.sendMsg, h.checkRegistration,
			func(ctx *tgf.Context) error {
				ctx.Log.Info("[HANDLER]: Display List Handler")

				c := EditListsHandlerContext{
					ShoppingLists: []types.ShoppingList{},
					Items:         []types.ShoppingListItem{},
				}

				if ctx.Journey.RawContext != nil {
					err := json.Unmarshal(ctx.Journey.RawContext, &c)
					if err != nil {
						ctx.Log.Error("Error unmarshaling context", "error", err)
						return fmt.Errorf("%w: %w", tgf.CouldNotExteactContextErr, err)
					}
				}

				chatID := ctx.GetChatID()
				message := ctx.GetMessage()

				data := ""
				if ctx.Update.CallbackQuery != nil {
					data = ctx.Update.CallbackQuery.Data
				}

				switch data {
				case "delete":
					ctx.Log.Info("deleting", "selected_list", c.SelectedList)

					if c.SelectedList == "" {
						return nil
					}

					err := h.deleteList(c.SelectedList)
					if err != nil {
						ctx.Log.Error("Error deleting list", "error", err)
						return fmt.Errorf("error deleting list from db: %w", err)
					}
					lists, err := h.getLists(chatID)
					if err != nil {
						return err
					}
					c.ShoppingLists = lists

				case "done":
					deleteMsg := tgbotapi.NewDeleteMessage(chatID, message.MessageID)
					_, err := h.botReq(deleteMsg)
					if err != nil {
						ctx.Log.Error("Error deleting inline keyboard", "error", err)
						return fmt.Errorf("error making bot request: %w", err)
					}
					ctx.Exit()
					return nil
				default:
					if data == c.SelectedList {
						c.SelectedList = ""
					} else if data != "" {
						c.SelectedList = data
					}

					lists, err := h.getLists(chatID)
					if err != nil {
						return err
					}

					if len(lists) < 1 {
						msg := tgbotapi.NewMessage(chatID, "There are no lists to chose from")
						_, err = h.sendMsg(msg)
						if err != nil {
							return err
						}
						ctx.Exit()
						return nil
					}
					c.ShoppingLists = lists
				}

				kbRows := h.buildKeyboard(c)

				if ctx.Update.Message != nil {
					msg := tgbotapi.NewMessage(chatID, "Please chose the list to edit")
					msg.ReplyMarkup = kbRows
					_, err := h.sendMsg(msg)
					if err != nil {
						ctx.Log.Error("Error sending message %w", err)
						return err
					}
				} else {
					msg := tgbotapi.NewEditMessageReplyMarkup(chatID, message.MessageID, kbRows)
					_, err := h.botReq(msg)
					if err != nil {
						ctx.Log.Error("Error sending bot request", "error", err)
						return fmt.Errorf("error making bot request: %w", err)
					}
				}

				ctx.Loop()
				return ctx.SetContexData(c)
			}),
	}
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

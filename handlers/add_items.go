package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/tgf"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddItemsHandler struct {
	sendMsg           func(c tgbotapi.Chattable) (tgbotapi.Message, error)
	botReq            func(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	getLists          func(chatID int64) ([]types.ShoppingList, error)
	addItems          func(listID string, itemText []string) error
	checkRegistration func(chatID int64) (bool, error)
}

func NewAddItemsHandler(
	msgSener func(c tgbotapi.Chattable) (tgbotapi.Message, error),
	botReq func(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error),
	getLists func(chatID int64) ([]types.ShoppingList, error),
	addItems func(listID string, itemText []string) error,
	checkRegistration func(chatID int64) (bool, error),
) *AddItemsHandler {
	return &AddItemsHandler{
		sendMsg:           msgSener,
		botReq:            botReq,
		getLists:          getLists,
		addItems:          addItems,
		checkRegistration: checkRegistration,
	}
}

type AddItemsHandlerContext struct {
	ShoppingListsMap map[string]types.ShoppingList
	ShoppingList     types.ShoppingList
	Items            []string
	ListEditable     bool
}

func (h *AddItemsHandler) GetHandlerJourney() []tgf.HandlerFunc {
	return []tgf.HandlerFunc{
		chatRegistered(h.sendMsg, h.checkRegistration),
		func(ctx *tgf.Context) error {
			ctx.Log.Info("[HANDLER]: Add Items Handler")

			chatID := ctx.GetChatID()
			lists, err := h.getLists(chatID)
			if err != nil {
				ctx.Log.Error("Error getting lists from db", "error", err)
				return err
			}

			if len(lists) < 1 {
				msg := tgbotapi.NewMessage(chatID, "There are no lists")
				_, err = h.sendMsg(msg)
				if err != nil {
					ctx.Log.Error("Error sending message", "error", err)
					return err
				}

				ctx.Exit()
				return nil
			}

			c := AddItemsHandlerContext{
				ShoppingListsMap: map[string]types.ShoppingList{},
			}

			for _, l := range lists {
				c.ShoppingListsMap[l.ID] = l
			}

			msg := tgbotapi.NewMessage(chatID, "Please chose the list to add items to")
			msg.ReplyMarkup = h.builldListsKeyboard(c)
			_, err = h.sendMsg(msg)
			if err != nil {
				ctx.Log.Error("Error sending message %w", err)
				return err
			}

			return ctx.SetContexData(c)
		},
		func(ctx *tgf.Context) error {
			ctx.Log.Info("[HANDLER]: Add Items Handler 2")
			chatID := ctx.GetChatID()
			listID := ctx.Update.CallbackQuery.Data

			var c AddItemsHandlerContext
			err := json.Unmarshal(ctx.Journey.RawContext, &c)
			if err != nil {
				ctx.Log.Error("Error unmarshaling context", "error", err)
				return fmt.Errorf("%w: %w", tgf.CouldNotExteactContextErr, err)
			}

			c.ShoppingList = c.ShoppingListsMap[listID]

			msg := tgbotapi.NewMessage(
				chatID,
				fmt.Sprintf("Adding to %s, start typing the items and type \"DONE\" when finished", c.ShoppingList.Title),
			)

			_, err = h.sendMsg(msg)
			if err != nil {
				ctx.Log.Error("Error sending message", "error", err)
				return err
			}

			return ctx.SetContexData(c)
		},
		func(ctx *tgf.Context) error {
			ctx.Log.Info("[HANDLER]: Add Items Handler 2")
			chatID := ctx.GetChatID()

			message := ctx.GetMessage()
			if message == nil {
				ctx.Exit()
				return nil
			}

			var c AddItemsHandlerContext
			err := json.Unmarshal(ctx.Journey.RawContext, &c)
			if err != nil {
				ctx.Log.Error("Error unmarshaling context", "error", err)
				return fmt.Errorf("%w: %w", tgf.CouldNotExteactContextErr, err)
			}

			if strings.ToUpper(message.Text) == "DONE" {
				if len(c.Items) == 0 {
					ctx.Exit()
					return nil
				}

				err := h.addItems(c.ShoppingList.ID, c.Items)
				if err != nil {
					ctx.Log.Error("Error adding items to db: %w", err)
					return fmt.Errorf("error inserting items into list %s, %w", c.ShoppingList.ID, err)
				}

				// TODO: smth is broken here with the commas
				textMessage := "Adding items: "
				for i, item := range c.Items {
					if i >= len(c.Items)-1 || i == 0 {
						textMessage = fmt.Sprintf("%s %s", textMessage, item)
						continue
					}
					textMessage = fmt.Sprintf("%s, %s", textMessage, item)
				}
				textMessage = fmt.Sprintf("%s to %s", textMessage, c.ShoppingList.Title)

				msg := tgbotapi.NewMessage(chatID, textMessage)
				_, err = h.sendMsg(msg)
				if err != nil {
					ctx.Log.Error("Error sending message", "error", err)
					return err
				}

				ctx.Exit()
				return nil
			}

			c.Items = append(c.Items, message.Text)

			ctx.Loop()
			return ctx.SetContexData(c)
		},
	}
}

func (h AddItemsHandler) builldListsKeyboard(c AddItemsHandlerContext) tgbotapi.InlineKeyboardMarkup {
	kbRows := [][]tgbotapi.InlineKeyboardButton{}
	for _, l := range c.ShoppingListsMap {

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s - %s", l.Title, l.StoreLocation), l.ID),
		)

		if c.ListEditable {
			row = append(row,
				tgbotapi.NewInlineKeyboardButtonData("❌", "del:"+l.ID),
				tgbotapi.NewInlineKeyboardButtonData("✏️", "change_name:"+l.ID),
			)

		}

		kbRows = append(
			kbRows,
			row,
		)

	}

	// TODO: Implement
	// kbRows = append(
	// 	kbRows,
	// 	tgbotapi.NewInlineKeyboardRow(
	// 		tgbotapi.NewInlineKeyboardButtonData("Close", "close"),
	// 		tgbotapi.NewInlineKeyboardButtonData("Edit", "edit"),
	// 	),
	// )

	return tgbotapi.NewInlineKeyboardMarkup(kbRows...)
}

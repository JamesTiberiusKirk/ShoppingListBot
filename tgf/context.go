package tgf

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(ctx *Context) error
type HandlerInterface interface {
	// GetHandlerJourney returns handler funcs journeys and weather or not the final element in the array is to be called endlessly
	GetHandlerJourney() []HandlerFunc
}

type Context struct {
	// Update - telegram-bot-api update object.
	Update tgbotapi.Update
	// Journey - TGF journey context object will be stored in the journey store.
	Journey *Journey
	// Log instance of Logger initialised with the update
	Log Logger

	// This object is meant to be for injecting dependencies
	// injected map[string]interface{}

	// handleError func(err error)
	// execHandler func()

	handlers       []HandlerFunc
	nextHasBeenSet bool
	skipTo         func(int)
	exit           func()
	changeJourney  func(string, int)
}

func (ctx *Context) SendMessage(text string) error {
	return nil
}

func (ctx *Context) SendRequest() error {
	return nil
}

// TODO: See if there is a way to either create a new keyboard or edit it if one exits
func (ctx *Context) UpsertInlineKeyboard(kb tgbotapi.InlineKeyboardMarkup) error {
	return nil
}

// GetChatID - returns chatid from the message with GetMessage(), returns 0 if no message found
func (ctx *Context) GetChatID() int64 {
	message := ctx.GetMessage()
	if message == nil {
		return 0
	}

	return message.Chat.ID
}

func (ctx *Context) SetContexData(data any) error {
	bytes, err := json.Marshal(data)
	ctx.Journey.RawContext = bytes
	return err
}

// GetMessage - gets message from update.message or update.CallbackQuery
func (ctx *Context) GetMessage() *tgbotapi.Message {
	if ctx.Update.Message != nil {
		return ctx.Update.Message
	} else if ctx.Update.CallbackQuery != nil {
		return ctx.Update.CallbackQuery.Message
	}

	return nil
}

func (ctx *Context) SkipTo(index int) {
	ctx.nextHasBeenSet = true
	ctx.skipTo(index)
}

func (ctx *Context) Loop() {
	ctx.nextHasBeenSet = true
}

func (ctx *Context) SetNextExec(index int) {
	if ctx.handlers[index] == nil {
		ctx.Exit()
		return
	}
	ctx.Journey.Next = index
}

func (ctx *Context) Exit() {
	ctx.exit()
}

func (ctx *Context) ChangeHandler(command string, index int) {
	ctx.changeJourney(command, index)
}

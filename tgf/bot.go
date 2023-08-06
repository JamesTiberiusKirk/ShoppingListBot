package tgf

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	CouldNotExteactContextErr = errors.New("could not extract context")
	UserErr                   = errors.New("user error")
	// TODO: Maybe think of making some error which would posibly just skip direclty to next handler?
)

func InitBotAPI(token string, telegramWebHookURL string, debug bool) (*tgbotapi.BotAPI, error) {
	// TODO: need to figure out how to setup webhooks propery
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug
	return bot, nil
}

type Bot struct {
	commands     tgbotapi.SetMyCommandsConfig
	journeyMap   map[string]HandlerInterface
	log          Logger
	bot          *tgbotapi.BotAPI
	journeyStore JourneyStore
}

func NewBot(bot *tgbotapi.BotAPI, commands tgbotapi.SetMyCommandsConfig, journeyMap map[string]HandlerInterface, log Logger, js JourneyStore) *Bot {
	return &Bot{
		bot:          bot,
		commands:     commands,
		journeyMap:   journeyMap,
		log:          log,
		journeyStore: js,
	}
}

func (b *Bot) SetLogger(logger Logger) {
	b.log = logger
}

func (b *Bot) StartBot(debug bool) error {
	if b.log == nil {
		b.log = NewDefaultLogger(debug)
	}

	b.log.Info("Authorized on account %s", b.bot.Self.UserName)

	_, err := b.bot.Request(b.commands)
	if err != nil {
		return err
	}

	botcfg := tgbotapi.NewUpdate(0)
	botcfg.Timeout = 60
	updates := b.bot.GetUpdatesChan(botcfg)
	for update := range updates {
		go b.handleUpdate(update)
	}

	return nil
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	b.log.SetUpdate(update)
	if update.Message == nil && update.CallbackQuery == nil {
		return
	}

	var message *tgbotapi.Message

	if update.Message != nil {
		message = update.Message
	}

	if update.CallbackQuery != nil {
		message = update.CallbackQuery.Message
	}

	if message == nil {
		return
	}

	ctx := &Context{
		Update:         update,
		nextHasBeenSet: false,
		Log:            b.log,
	}
	chatID := ctx.GetChatID()

	if message.IsCommand() {
		c := message.Command()
		ctx.Journey = &Journey{}
		ctx.Journey.Command = c
		ctx.Journey.Next = 0
		ctx.Journey.TelegramChatID = chatID

	} else {
		j, err := b.journeyStore.GetJourneyByChatID(chatID)
		if err != nil {
			b.log.Error("[JOURNEY ERROR]: trying to access handler journey store: %w", err.Error())
			return
		}

		if j == nil {
			b.log.Error("[JOURNEY ERROR]: journey nil")
			return
		}

		ctx.Journey = j
	}

	b.createCallbacks(ctx)

	err := b.getJourney(ctx)
	if err != nil {
		b.log.Error("[HANDLER ERROR]: getting handler: %w", err)
		return
	}

	err = b.execHandler(ctx)
	if err != nil {
		b.log.Error("[HANDLER ERROR]: error executing handler: %w", err)
		return
	}

	// Joruney exit
	if ctx.Journey == nil {
		return
	}

	if len(ctx.handlers) > ctx.Journey.Next {
		if !ctx.nextHasBeenSet {
			if ctx.handlers == nil && ctx.handlers[ctx.Journey.Next+1] == nil {
				ctx.Exit()
				return
			}
			ctx.nextHasBeenSet = true
			ctx.Journey.Next += 1
		}
		_, err := b.journeyStore.UpsertJourneyByTelegeramChatID(chatID, *ctx.Journey)
		if err != nil {
			b.log.Error("[JOURNEY ERROR]: trying to upsert handler journey DB: %w", err)
			return
		}
		return
	}

	// if this is the last in the journey, cleanup
	b.log.Debug("[SCHEDULER]: cleaning up index: %s", ctx.Journey.Next)
	err = b.journeyStore.CleanupChatJourney(chatID)
	if err != nil {
		b.log.Error("[DB ERROR]: trying to cleanup handler journey DB: %w", err)
		return
	}
}

func (b *Bot) createCallbacks(ctx *Context) {
	ctx.skipTo = func(i int) {
		if ctx.handlers[i] == nil {
			b.log.Error("[SKIP TO ERROR]: invalid index %d", i)
			ctx.Exit()
			return
		}

		ctx.Journey.Next = i
		b.handleHandlerError(ctx, b.execHandler(ctx))
		ctx.nextHasBeenSet = false
	}

	ctx.exit = func() {
		b.log.Debug("[SCHEDULER]: cleaning up index: %s", ctx.Journey.Next)
		err := b.journeyStore.CleanupChatJourney(ctx.GetChatID())
		if err != nil {
			b.log.Error("[DB ERROR]: trying to cleanup handler journey DB: %w", err)
			return
		}

		ctx.Journey = nil
	}

	ctx.changeJourney = func(command string, i int) {
		if command != ctx.Journey.Command {
			b.getJourney(ctx)
		}

		if ctx.handlers[i] == nil {
			b.log.Error("[CHANGE HANDLER ERROR]: invalid index %d", i)
			return
		}

		ctx.Journey.Next = i
	}
}

func (b *Bot) handleHandlerError(ctx *Context, err error) {
	if err != nil {
		return
	}

	b.log.Debug("[SCHEDULER]: cleaning up index: %s", ctx.Journey.Next)
	err = b.journeyStore.CleanupChatJourney(ctx.GetChatID())
	if err != nil {
		b.log.Error("[DB ERROR]: trying to cleanup handler journey DB: %w", err)
		return
	}
}

func (b *Bot) getJourney(ctx *Context) error {
	handlerJourney, ok := b.journeyMap[ctx.Journey.Command]
	if !ok {
		return cmdNotImplementedErr
	}
	journey := handlerJourney.GetHandlerJourney()
	if journey[ctx.Journey.Next] == nil {
		return errors.New("error trying to access handler journey which is nil")
	}

	ctx.handlers = journey
	return nil
}

var (
	cmdNotImplementedErr = errors.New("command not implemented")
)

func (b *Bot) execHandler(ctx *Context) error {
	chatID := ctx.GetChatID()

	// journey exit
	if ctx.Journey == nil {
		return nil
	}

	err := ctx.handlers[ctx.Journey.Next](ctx)
	if err != nil {
		if err != UserErr {
			b.log.Error("[HANDLER ERROR]: ", "chatID", chatID, "error", err)
			msg := tgbotapi.NewMessage(chatID, "Sorry, internal server error")
			if _, err := b.bot.Send(msg); err != nil {
				return fmt.Errorf("[HANDLER ERROR]: failed to send message %w", err)
			}
		}

		b.log.Debug("[SCHEDULER]: cleaning up", "chatID", chatID, "index", ctx.Journey.Next)
		errCleanup := b.journeyStore.CleanupChatJourney(chatID)
		if errCleanup != nil {
			return fmt.Errorf("[HANDLER ERROR]: trying to cleanup handler journey %w", errCleanup)
		}
		return fmt.Errorf("[HANDLER ERROR]: %w", err)
	}

	return nil
}

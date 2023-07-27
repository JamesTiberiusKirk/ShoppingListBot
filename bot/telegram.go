package bot

import (
	"github.com/JamesTiberiusKirk/ShoppingListsBot/db"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/inconshreveable/log15"
)

func StartBot(token string, telegramWebHookURL string, debug bool, db *db.DB) error {
	var bot *tgbotapi.BotAPI
	var err error
	var updates tgbotapi.UpdatesChannel

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	log.Info("Authorized on account " + bot.Self.UserName)
	bot.Debug = debug
	botcfg := tgbotapi.NewUpdate(0)
	botcfg.Timeout = 60

	commands := handlers.GetHandlerCommandList()
	_, err = bot.Request(commands)

	updates = bot.GetUpdatesChan(botcfg)

	// TODO: need to make sure this never gets manipulated
	journeyMap := handlers.NewHandlerJounreyMap(bot, db)

	for update := range updates {
		go handleUpdate(update, journeyMap, db, bot)
	}

	return nil
}

func handleUpdate(update tgbotapi.Update, journeyMap map[string]handlers.HandlerInterface, db *db.DB, bot *tgbotapi.BotAPI) {
	if update.Message == nil && update.CallbackQuery == nil {
		log.Info("skipping")
		return
	}

	var message *tgbotapi.Message

	if update.Message != nil {
		command := update.Message.Command()
		chatID := update.Message.Chat.ID
		message = update.Message

		log.Info("[HANDLER]:", "ChatID", chatID, "COMMAND", command, "TEXT", message.Text)
	}

	if update.CallbackQuery != nil {
		command := update.CallbackQuery.Message.Command()
		chatID := update.CallbackQuery.Message.Chat.ID
		message = update.CallbackQuery.Message

		log.Info("[CALLBACK QUERY HANDLER]:", "ChatID", chatID, "COMMAND", command, "TEXT", message.Text)
	}

	if message != nil {
		// TODO: cleanup this mess
		// TODO: NEED TO figure out a way to manage groups in here
		// So far a workaround I have found is to set the bot as admin in the groupchat
		command := ""
		index := 0
		var previousContext []byte

		if message.IsCommand() {
			command = message.Command()
		} else {
			c, err := db.GetJourneyByChat(message.Chat.ID)
			if err != nil {
				// HANDLE DB ERROR
				log.Error("[HANDLER ERROR]: trying to access handler journey DB", "chatID", message.Chat.ID, "error", err.Error())
				return
			}

			command = c.Command
			index = c.Next
			previousContext = c.RawContext
		}

		chatID := message.Chat.ID

		handlerJourney, ok := journeyMap[command]
		if !ok {
			// command not found, for now just ignore it for now
			return
		}

		journey, infinite := handlerJourney.GetHandlerJourney()
		if journey[index] == nil {
			log.Error("[HANDLER ERROR]: trying to access handler journey which is nil", "chatID", message.Chat.ID)
			return
		}

		nextContext, err := journey[index](previousContext, update)
		if err != nil {
			if err != handlers.JourneryExitErr && err != handlers.UserErr {
				log.Error("[HANDLER ERROR]:", "chatID", message.Chat.ID, "error", err.Error())
				msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, internal server error")
				if _, err := bot.Send(msg); err != nil {
					log.Error("[HANDLER ERROR]: failed to send message", "chatID", message.Chat.ID, "error", err)
					return
				}
			}

			log.Info("[SCHEDULER]: cleaning up", "chatID", chatID, "index", index)
			err := db.CleanupChatJourney(chatID)
			if err != nil {
				log.Error("[HANDLER ERROR]: trying to cleanup handler journey DB", "chatID", chatID, "error", err)
				return
			}
			return
		}

		if len(journey)-1 > index {
			_, err := db.UpsertJourneyByTelegeramChatID(chatID, command, index+1, nextContext)
			if err != nil {
				// HANDLE DB ERROR
				log.Error("[HANDLER ERROR]: trying to upsert handler journey DB", "chatID", chatID, "error", err)
				return
			}

			log.Info("[SCHEDULER]: next", "chatID", chatID, "index", index+1)
			return
		}

		if infinite {
			log.Info("[SCHEDULER]: infinite %d, %d", chatID, index)
			_, err := db.UpsertJourneyByTelegeramChatID(chatID, command, index, nextContext)
			if err != nil {
				// HANDLE DB ERROR
				log.Info("[INFINITE HANDLER ERROR]: trying to upsert handler journey DB", "chatID", chatID, "error", err)
				return
			}
			return
		}

		// if this is the last in the journey, cleanup
		log.Info("[SCHEDULER]: cleaning up", "chatID", chatID, "index", index)
		err = db.CleanupChatJourney(chatID)
		if err != nil {
			// HANDLE DB ERROR
			log.Info("[DB ERROR]: trying to cleanup handler journey DB", "chatID", chatID, "error", err)
			return
		}
	}
}

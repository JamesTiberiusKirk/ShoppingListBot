package bot

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/db"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(token string, telegramWebHookURL string, debug bool, db *db.DB) error {
	var bot *tgbotapi.BotAPI
	var err error
	var updates tgbotapi.UpdatesChannel

	if telegramWebHookURL != "" {
		bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			return err
		}

		bot.Debug = debug
		log.Printf("Authorized on account %s", bot.Self.UserName)

		wh, err := tgbotapi.NewWebhook(fmt.Sprintf("%s/%s", telegramWebHookURL, bot.Token))
		if err != nil {
			return err
		}

		_, err = bot.Request(wh)
		if err != nil {
			return err
		}

		info, err := bot.GetWebhookInfo()
		if err != nil {
			return err
		}

		if info.LastErrorDate != 0 {
			log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}

		commands := handlers.GetHandlerCommandList()
		_, err = bot.Request(commands)

		updates = bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServe("0.0.0.0:443", nil)

	} else {
		bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			return err
		}

		log.Printf("Authorized on account %s", bot.Self.UserName)
		bot.Debug = debug
		botcfg := tgbotapi.NewUpdate(0)
		botcfg.Timeout = 60

		commands := handlers.GetHandlerCommandList()
		_, err = bot.Request(commands)

		updates = bot.GetUpdatesChan(botcfg)
	}

	jouneyMap := handlers.NewHandlerJounreyMap(bot, db)

	for update := range updates {
		// TODO: wrap everything in a gorutine? dont forget to use the apropriate map type for gorutines
		if update.Message == nil && update.CallbackQuery == nil {
			log.Print("skipping")
			continue
		}

		var message *tgbotapi.Message

		if update.Message != nil {
			command := update.Message.Command()
			chatID := update.Message.Chat.ID
			message = update.Message

			log.Printf("[HANDLER]: ChatID: %d COMMAND: %s TEXT: %s", chatID, command, message.Text)
		}

		if update.CallbackQuery != nil {
			command := update.CallbackQuery.Message.Command()
			chatID := update.CallbackQuery.Message.Chat.ID
			message = update.CallbackQuery.Message

			log.Printf("[CALLBACK QUERY HANDLER]: ChatID: %d COMMAND: %s TEXT: %s", chatID, command, message.Text)
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
					log.Printf("[HANDLER ERROR]: chatID %d, trying to ccess handler journey DB error: %s", message.Chat.ID, err.Error())
					continue
				}

				command = c.Command
				index = c.Next
				previousContext = c.RawContext
			}

			chatID := message.Chat.ID

			handlerJourney, ok := jouneyMap[command]
			if !ok {
				// command not found, for now just ignore it
				log.Printf("[HANDLER ERROR]: chatID %d, trying to ccess handler journey which is nil", message.Chat.ID)
				continue
			}
			journey, infinite := handlerJourney.GetHandlerJourney()
			if journey[index] == nil {
				log.Printf("[HANDLER ERROR]: chatID %d, trying to ccess handler journey which is nil", message.Chat.ID)
				continue
			}

			nextContext, err := journey[index](previousContext, update)
			if err != nil {
				if err != handlers.JourneryExitErr && err != handlers.UserErr {
					log.Printf("[HANDLER ERROR]: ChatID: %d, %s", message.Chat.ID, err.Error())
					msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, internal server error")
					bot.Send(msg)
				}

				log.Printf("cleaning up %d, %d", chatID, index)
				err := db.CleanupChatJourney(chatID)
				if err != nil {
					// HANDLE DB ERROR
					log.Printf("[HANDLER ERROR]: chatID %d, trying to cleanup handler journey DB error: %s", chatID, err.Error())
					continue
				}
				continue
			}

			if len(journey)-1 > index {
				_, err := db.UpsertJourneyByTelegeramChatID(chatID, command, index+1, nextContext)
				if err != nil {
					// HANDLE DB ERROR
					log.Printf("[HANDLER ERROR]: chatID %d, trying to upsert handler journey DB error: %s", chatID, err.Error())
					continue
				}

				log.Printf("next %d, %d", chatID, index)
				continue
			}

			if infinite {
				log.Printf("infinite %d, %d", chatID, index)
				_, err := db.UpsertJourneyByTelegeramChatID(chatID, command, index, nextContext)
				if err != nil {
					// HANDLE DB ERROR
					log.Printf("[INFINITE HANDLER ERROR]: chatID %d, trying to upsert handler journey DB error: %s", chatID, err.Error())
					continue
				}
				continue
			}

			// if this is the last in the journey, cleanup
			log.Printf("cleaning up %d, %d", chatID, index)
			err = db.CleanupChatJourney(chatID)
			if err != nil {
				// HANDLE DB ERROR
				log.Printf("[HANDLER ERROR]: chatID %d, trying to cleanup handler journey DB error: %s", chatID, err.Error())
				continue
			}
		}
	}

	return nil
}

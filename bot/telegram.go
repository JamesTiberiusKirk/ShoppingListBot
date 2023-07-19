package bot

import (
	"log"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/clients"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(token string, debug bool, db *clients.DB) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	var botcfg tgbotapi.UpdateConfig

	switch debug {
	case true:
		// This thorws a lot more console logs
		bot.Debug = false
		log.Printf("Authorized on account %s", bot.Self.UserName)
		botcfg = tgbotapi.NewUpdate(0)
		botcfg.Timeout = 60
	case false:
		panic("Production mode not implemented yet")
	}

	updates := bot.GetUpdatesChan(botcfg)

	jouneyMap := handlers.NewHandlerJounreyMap(bot, db)

	commands := handlers.GetHandlerCommandList()
	_, err = bot.Request(commands)

	// TODO: This needs to memcached or redis
	// This way I can add timed cleanup there in the form of expiry time (not sure if i can in memcached)
	contexHandlerTracker := map[int64]handlers.JourneyTracker{}

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
			log.Printf("[HANDLER]: ChatID: %d COMMAND: %s", chatID, command)
			message = update.Message

		}

		if update.CallbackQuery != nil {
			command := update.CallbackQuery.Message.Command()
			chatID := update.CallbackQuery.Message.Chat.ID
			log.Printf("[CALLBACK COMMAND]: %d, COMMAND: %s, DATA: %s", chatID,
				command, update.CallbackQuery.Data)
			message = update.CallbackQuery.Message
		}

		if message != nil {
			command := ""
			index := 0
			previous := []tgbotapi.Update{}

			if message.IsCommand() {
				command = message.Command()
			} else {
				c, ok := contexHandlerTracker[message.Chat.ID]
				if !ok {
					continue
				}
				command = c.Command
				index = c.Next
				previous = c.PastUpdates
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

			err := journey[index](update, previous)
			if err != nil {
				if err != handlers.JourneryExitErr {
					log.Printf("[HANDLER ERROR]: ChatID: %d, %s", message.Chat.ID, err.Error())
					msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, internal server error")
					bot.Send(msg)
				}

				log.Printf("cleaning up %d, %d", chatID, index)
				delete(contexHandlerTracker, chatID)
				continue
			}

			if len(journey)-1 > index {
				contexHandlerTracker[chatID] = handlers.JourneyTracker{
					Next:        index + 1,
					Command:     command,
					PastUpdates: append(previous, update),
				}
				log.Printf("next %d, %d", chatID, index)
				continue
			}

			if infinite {
				log.Printf("infinite %d, %d", chatID, index)
				contexHandlerTracker[chatID] = handlers.JourneyTracker{
					Next:        index,
					Command:     command,
					PastUpdates: append(previous, update),
				}
				continue
			}

			// if this is the last in the journey, cleanup
			log.Printf("cleaning up %d, %d", chatID, index)
			delete(contexHandlerTracker, chatID)
		}
	}

	return nil
}

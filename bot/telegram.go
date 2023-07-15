package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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

	handlerMap := handlers.NewHandlerMap(bot, db)
	callbackMap := handlers.NewCallbackMap(bot, db)
	contextualReplyMap := handlers.NewReplyCallbackMap(bot, db)

	commands := handlers.GetHandlerCommandList()
	_, err = bot.Request(commands)

	contextualReplyHandlerLookup := map[int64]string{}

	for update := range updates {
		if update.Message != nil {
			if !update.Message.IsCommand() {
				// TODO: Need to figure out here how to handle messages which came after specific commands
				// Might need todo smth with chatID
				// E.G. assuming some form of in memroy db as a map
				// Use chatID as key and value would be the handler command
				// Then if message is not a commend
				//	Lookup the chatID in the map
				//	if not found
				//		continue
				//	getContextualReplyHandler with the message command
				//	NOTE: Will need to make a cleanup function

				// NOTE: The above approach works well with just one contextual reply, but I need to figure out how to chain multiple

				command, ok := contextualReplyHandlerLookup[update.Message.Chat.ID]
				if !ok {
					continue
				}

				contextualReplyHandler, ok := contextualReplyMap[command]
				if !ok {
					continue
				}

				exeList := strings.Split(command, ":")
				if exeList[1] != "~" {
					exeNum, err := strconv.Atoi(exeList[1])
					if err != nil {
						log.Printf("[CONTEXTUAL CALLBACK ERROR]: ChatID: %d, %s", update.Message.Chat.ID, err.Error())
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, internal server error")
						bot.Send(msg)
						continue
					}

					// Check for a next one, else if check for an infinitly looping one
					newCommand := fmt.Sprintf("%s:%d", exeList[0], exeNum+1)
					_, ok = contextualReplyMap[newCommand]
					if !ok {
						newCommand = fmt.Sprintf("%s:~", exeList[0])
						_, ok := contextualReplyMap[newCommand]
						if !ok {
							delete(contextualReplyHandlerLookup, update.Message.Chat.ID)
							continue
						}
					}

					contextualReplyHandlerLookup[update.Message.Chat.ID] = newCommand
				}

				err = contextualReplyHandler(update)
				if err != nil {
					log.Printf("[CONTEXTUAL CALLBACK ERROR]: ChatID: %d, %s", update.Message.Chat.ID, err.Error())
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, internal server error")
					bot.Send(msg)
					continue
				}

				continue
			}

			log.Printf("[CALLBACK]: ChatID: %d, DATA: %s", update.Message.Chat.ID, update.Message.Text)

			handler, ok := handlerMap[update.Message.Command()]
			if !ok {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No such command")
				bot.Send(msg)
				continue
			}

			contextualReplyHandlerLookup[update.Message.Chat.ID] = fmt.Sprintf("%s:0", update.Message.Command())
			err := handler(update)
			if err != nil {
				log.Printf("[HANDLER ERROR]: ChatID: %d, %s", update.Message.Chat.ID, err.Error())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, internal server error")
				bot.Send(msg)
				continue
			}

		} else if update.CallbackQuery != nil {
			log.Printf("[CALLBACK]: %s, DATA: %s", update.CallbackQuery.ID, update.CallbackQuery.Data)

			if !update.CallbackQuery.Message.IsCommand() {
				continue
			}

			callback, ok := callbackMap[update.CallbackQuery.Message.Command()]
			if !ok {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "404")
				bot.Send(msg)
				continue
			}
			err := callback(update)
			if err != nil {
				log.Printf("[CALLBACK ERROR]: %s", err.Error())
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Sorry, internal server error")
				bot.Send(msg)
				continue
			}
		}
	}
	return nil
}

package bot

import (
	"log"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(token string, debug bool) error {
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

	handlerMap := handlers.NewHandlerMap(bot)
	callbackMap := handlers.NewCallbackMap(bot)

	commands := handlers.GetHandlerCommandList()
	_, err = bot.Request(commands)

	for update := range updates {
		if update.Message == nil {

			if update.CallbackQuery != nil {
				log.Printf("[CALLBACK]: %s, DATA: %s", update.CallbackQuery.ID, update.CallbackQuery.Data)
				log.Printf("%s", update.CallbackQuery.Message.Text)

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

			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		handler, ok := handlerMap[update.Message.Command()]
		if !ok {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "404")
			bot.Send(msg)
			continue
		}
		err := handler(update)
		if err != nil {
			log.Printf("[HANDLER ERROR]: %s", err.Error())
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, internal server error")
			bot.Send(msg)
			continue
		}
	}
	return nil
}

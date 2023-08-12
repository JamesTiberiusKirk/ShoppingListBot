package db

import (
	"github.com/JamesTiberiusKirk/ShoppingListsBot/tgf"
)

type DBJourneyStore struct {
	db *DB
}

func NewDBJourneyStore(db *DB) *DBJourneyStore {
	return &DBJourneyStore{
		db: db,
	}
}

func (js *DBJourneyStore) GetJourneyByChatID(chatID int64) (*tgf.Journey, error) {
	j, err := js.db.GetJourneyByChat(chatID)
	if err != nil {
		return nil, err
	}

	return &tgf.Journey{
		ID:             j.ID,
		TelegramChatID: j.TelegramChatID,
		ChatID:         j.ChatID,
		Command:        j.Command,
		Next:           j.Next,
		RawContext:     j.RawContext,
	}, nil
}
func (js *DBJourneyStore) CleanupChatJourney(chatID int64) error {
	return js.db.CleanupChatJourney(chatID)
}

func (js *DBJourneyStore) UpsertJourneyByTelegeramChatID(chatID int64, upsert tgf.Journey) (*tgf.Journey, error) {
	upsertedJourney, err := js.db.UpsertJourneyByTelegeramChatID(chatID, upsert.Command, upsert.Next, upsert.RawContext)
	if err != nil {
		return nil, err
	}

	return &tgf.Journey{
		ID:             upsertedJourney.ID,
		TelegramChatID: upsertedJourney.TelegramChatID,
		ChatID:         upsertedJourney.ChatID,
		Command:        upsertedJourney.Command,
		Next:           upsertedJourney.Next,
		RawContext:     upsertedJourney.RawContext,
	}, nil
}

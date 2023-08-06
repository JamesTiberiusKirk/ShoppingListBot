package db

import (
	"database/sql"
	"fmt"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	log "github.com/inconshreveable/log15"
)

func (d *DB) UpsertJourneyByTelegeramChatID(chatID int64, command string, next int, context []byte) (*types.Journey, error) {
	qName := "upsert_journey_by_telegram_chat_id"
	log.Info("[DB]: upserting journey", "query_name", qName, "chatID", chatID, "command", command, "next", next)
	query, ok := d.queries[qName]
	if !ok {
		return nil, fmt.Errorf("query missing %s", qName)
	}

	namedExecParams := map[string]interface{}{
		"telegram_chat_id": chatID,
		"command":          command,
		"next":             next,
		"context":          context,
	}

	rows, err := d.db.NamedQuery(query.Query, namedExecParams)
	if err != nil {
		return nil, fmt.Errorf("error upserting into chats_journey %d, err: %w", chatID, err)
	}

	var updatedJourney types.Journey
	for rows.Next() {
		err = rows.StructScan(&updatedJourney)
		if err != nil {
			return nil, fmt.Errorf("error getting updated chats_journey %d, err: %w", chatID, err)
		}
	}

	return &updatedJourney, nil
}

func (d *DB) GetJourneyByChat(chatID int64) (*types.Journey, error) {
	qName := "get_chat_journey_by_telegram_chat_id"
	log.Info("[DB]: get journey by chat", "query_name", qName, "chatID", chatID)
	query, ok := d.queries[qName]
	if !ok {
		return nil, fmt.Errorf("query missing %s", qName)
	}

	var journey types.Journey
	err := d.db.QueryRowx(query.Query, chatID).StructScan(&journey)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error finding chat journey by id: %d err: %w", chatID, err)
	}

	return &journey, nil
}

func (d *DB) CleanupChatJourney(chatID int64) error {
	qName := "cleanup_chat_journies_by_telegram_chat_id"
	log.Info("[DB]: cleaning up journey", "query_name", qName, "chatID", chatID)
	query, ok := d.queries[qName]
	if !ok {
		return fmt.Errorf("query missing %s", qName)
	}

	_, err := d.db.Exec(query.Query, chatID)
	if err != nil {
		return fmt.Errorf("error cleaning up journey %d err: %w", chatID, err)
	}

	return nil
}

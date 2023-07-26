package db

import (
	"database/sql"
	"fmt"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	log "github.com/inconshreveable/log15"
)

func (d *DB) AddNewChat(chatID int64) error {
	addChatSQL, ok := d.queries["add_chat"]
	if !ok {
		return fmt.Errorf("query missing add_chat")
	}

	_, err := d.db.Exec(addChatSQL.Query, chatID)
	if err != nil {
		return fmt.Errorf("error inserting into chats %d", chatID)
	}

	return nil
}

func (d *DB) CheckIfChatExists(chatID int64) (bool, error) {
	log.Info("[DB] quering chats table for chat", "chatID", chatID)

	query, ok := d.queries["get_chat"]
	if !ok {
		return false, fmt.Errorf("query missing get_chat")
	}

	var chat types.Chat
	err := d.db.QueryRowx(query.Query, chatID).StructScan(&chat)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error quering chats table: %w", err)
	}

	if chatID != chat.TelegramChatID {
		return false, nil
	}

	return true, nil
}

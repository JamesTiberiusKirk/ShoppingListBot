package clients

import (
	"fmt"
	"log"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql"

	_ "github.com/lib/pq"
)

type DB struct {
	db      *sqlx.DB
	schema  goyesql.Queries
	queries goyesql.Queries
}

func NewDBClient(dbUrl string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	schema := goyesql.MustParseFile("./sql/schema.sql")
	queries := goyesql.MustParseFile("./sql/queries.sql")

	sq, ok := schema["schema"]
	if !ok {
		log.Print("query not found")
	}

	_, err = db.Exec(sq.Query)
	if err != nil {
		return nil, err
	}

	return &DB{
		db:      db,
		schema:  schema,
		queries: queries,
	}, nil
}

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
	selectChat, ok := d.queries["get_chat"]
	if !ok {
		return false, fmt.Errorf("query missing get_chat")
	}

	chats := []types.Chat{}
	err := d.db.Select(&chats, selectChat.Query, chatID)
	if err != nil {
		return false, fmt.Errorf("error inserting into chats %d, err: %w", chatID, err)
	}

	if len(chats) >= 1 {
		log.Printf("%+v", chats)
		return false, fmt.Errorf("multiple chats found")
	}

	if len(chats) != 0 && chats[0].ChatID != chatID {
		return true, nil
	}

	return false, nil
}

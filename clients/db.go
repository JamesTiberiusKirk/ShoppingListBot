package clients

import (
	"fmt"
	"log"
	"time"

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

	return &DB{
		db:      db,
		schema:  schema,
		queries: queries,
	}, nil
}

func (d *DB) ApplySchema() error {
	sq, ok := d.schema["schema"]
	if !ok {
		log.Print("schema not found")
		return fmt.Errorf("schemanot not found")
	}

	_, err := d.db.Exec(sq.Query)
	if err != nil {
		return err
	}

	return nil
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

	if len(chats) != 0 && chats[0].TelegramChatID != chatID {
		return true, nil
	}

	return false, nil
}

func (d *DB) NewShoppingList(chatID int64, title string, storeLoc string, dueDate *time.Time) error {
	log.Printf("[DB] inserting to shopping_lists: %+v, %+v, %+v, %+v", chatID, title, storeLoc, dueDate)

	addListQuery, ok := d.queries["add_list"]
	if !ok {
		return fmt.Errorf("query missing add_list")
	}

	_, err := d.db.Exec(addListQuery.Query, chatID, title, storeLoc, dueDate)
	if err != nil {
		return fmt.Errorf("error inserting into shopping list table: %w", err)
	}

	return nil
}

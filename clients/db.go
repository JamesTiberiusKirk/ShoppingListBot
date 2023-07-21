package clients

import (
	"database/sql"
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
		return fmt.Errorf("error inserting into shopping_lists table: %w", err)
	}

	return nil
}

func (d *DB) GetListsByID(id int64) (types.ShoppingList, error) {
	log.Printf("[DB] quering shopping_lists table by id: %+v", id)

	addListQuery, ok := d.queries["get_list_by_id"]
	if !ok {
		return types.ShoppingList{}, fmt.Errorf("query missing get_list_by_id")
	}

	var shoppingLists types.ShoppingList
	err := d.db.Select(&shoppingLists, addListQuery.Query, id)
	if err != nil {
		return types.ShoppingList{}, fmt.Errorf("error quering shopping_lists table: %w", err)
	}

	return shoppingLists, nil
}

func (d *DB) GetListsByChat(chatID int64) ([]types.ShoppingList, error) {
	log.Printf("[DB] quering shopping_lists table for chat: %+v", chatID)

	addListQuery, ok := d.queries["get_lists"]
	if !ok {
		return nil, fmt.Errorf("query missing get_lists")
	}

	var shoppingLists []types.ShoppingList
	err := d.db.Select(&shoppingLists, addListQuery.Query, chatID)
	if err != nil {
		return nil, fmt.Errorf("error quering shopping_lists table: %w", err)
	}

	return shoppingLists, nil
}

func (d *DB) AddItemsToList(listID string, itemsText []string) error {
	log.Printf("[DB] inserting to shopping_list_items: %+v, %+v", listID, itemsText)

	query, ok := d.queries["add_items"]
	if !ok {
		return fmt.Errorf("query missing add_items")
	}

	batchInsert := []map[string]interface{}{}

	for _, t := range itemsText {
		batchInsert = append(batchInsert, map[string]interface{}{
			"list_id": listID,
			"text":    t,
		})
	}

	_, err := d.db.NamedExec(query.Query, batchInsert)
	if err != nil {
		return fmt.Errorf("error inserting into shopping_list_items table: %w", err)
	}

	return nil
}

func (d *DB) GetItemsByList(listID string) ([]types.ShoppingListItem, error) {
	log.Printf("[DB] quering shopping_list_items table for chat: %+v", listID)

	query, ok := d.queries["get_items_in_list"]
	if !ok {
		return nil, fmt.Errorf("query missing get_items_in_list")
	}

	var shoppingListItems []types.ShoppingListItem
	err := d.db.Select(&shoppingListItems, query.Query, listID)
	if err != nil {
		return nil, fmt.Errorf("error quering shopping_list_items table: %w", err)
	}

	return shoppingListItems, nil
}

func (d *DB) ToggleItemPurchase(itemID string) error {
	log.Printf("[DB] toggling item purchase: %+v", itemID)

	query, ok := d.queries["toggle_item_purchase"]
	if !ok {
		return fmt.Errorf("query missing toggle_item_purchase")
	}

	_, err := d.db.Exec(query.Query, itemID)
	if err != nil {
		return fmt.Errorf("error updating toggle_item_purchase table: %w", err)
	}

	return nil
}

func (d *DB) CheckRegistration(chatID int64) (bool, error) {
	log.Printf("[DB] quering chats table for chat: %+v", chatID)

	query, ok := d.queries["get_chat"]
	if !ok {
		return false, fmt.Errorf("query missing get_chat")
	}

	var chat types.Chat
	err := d.db.QueryRow(query.Query, chatID).Scan(&chat)
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

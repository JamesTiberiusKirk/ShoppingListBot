package db

import (
	"fmt"
	"time"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	log "github.com/inconshreveable/log15"
)

func (d *DB) NewShoppingList(chatID int64, title string, storeLoc string, dueDate *time.Time) error {
	qName := "add_list"
	log.Info("[DB]: inserting to shopping_lists", "query_name", qName, "chatID", chatID, "title", title, "store_loc", storeLoc, "due_date", dueDate)
	addListQuery, ok := d.queries[qName]
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
	qName := "get_list_by_id"
	log.Info("[DB]: inserting to shopping_lists", "query_name", qName, "id", id)

	addListQuery, ok := d.queries[qName]
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
	qName := "get_lists"
	log.Info("[DB] quering shopping_lists table for chat", "chatID", chatID)

	addListQuery, ok := d.queries[qName]
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

func (d *DB) DeleteListByID(id string) error {
	qName := "delete_list_by_id"
	log.Info("[DB]: delete to shopping_lists", "query_name", qName, "id", id)

	addListQuery, ok := d.queries[qName]
	if !ok {
		return fmt.Errorf("query missing get_list_by_id")
	}

	_, err := d.db.Query(addListQuery.Query, id)
	if err != nil {
		return fmt.Errorf("error deleting from shopping_lists table: %w", err)
	}

	return nil
}

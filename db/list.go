package db

import (
	"fmt"
	"log"
	"time"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
)

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

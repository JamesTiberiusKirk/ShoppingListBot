package db

import (
	"fmt"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
	log "github.com/inconshreveable/log15"
)

func (d *DB) AddItemsToList(listID string, itemsText []string) error {
	log.Info("[DB] inserting to shopping_list_items", "listID", listID, "itemsText", itemsText)

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
	log.Info("[DB] quering shopping_list_items table for chat:", "listID", listID)

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
	log.Info("[DB] toggling item purchase:", "itemID", itemID)

	query, ok := d.queries["toggle_item_purchase"]
	if !ok {
		return fmt.Errorf("query missing toggle_item_purchase")
	}

	_, err := d.db.Exec(query.Query, itemID)
	if err != nil {
		return fmt.Errorf("error updating items table: %w", err)
	}

	return nil
}

func (d *DB) DeleteItem(itemID string) error {
	log.Info("[DB] deleting item:", "itemID", itemID)

	qName := "delete_item"
	query, ok := d.queries[qName]
	if !ok {
		return fmt.Errorf("query missing %s", qName)
	}

	_, err := d.db.Exec(query.Query, itemID)
	if err != nil {
		return fmt.Errorf("error deleting from items table: %w", err)
	}

	return nil
}

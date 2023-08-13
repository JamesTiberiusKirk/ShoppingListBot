package db

import (
	"fmt"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/types"
)

func (d *DB) AddItemsToList(listID string, itemsText []string) error {
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

	_, err := d.DB.NamedExec(query.Query, batchInsert)
	if err != nil {
		return fmt.Errorf("error inserting into shopping_list_items table: %w", err)
	}

	return nil
}

func (d *DB) GetItemsByList(listID string, showPurchased bool) ([]types.ShoppingListItem, error) {
	qName := "get_unpurchased_items_in_list"
	if showPurchased {
		qName = "get_items_in_list"
	}

	query, ok := d.queries[qName]
	if !ok {
		return nil, fmt.Errorf("query missing %s", qName)
	}

	var shoppingListItems []types.ShoppingListItem
	err := d.DB.Select(&shoppingListItems, query.Query, listID)
	if err != nil {
		return nil, fmt.Errorf("error qName: %s, quering shopping_list_items table: %w", qName, err)
	}

	return shoppingListItems, nil
}

func (d *DB) ToggleItemPurchase(itemID string) error {
	query, ok := d.queries["toggle_item_purchase"]
	if !ok {
		return fmt.Errorf("query missing toggle_item_purchase")
	}

	_, err := d.DB.Exec(query.Query, itemID)
	if err != nil {
		return fmt.Errorf("error updating items table: %w", err)
	}

	return nil
}

func (d *DB) DeleteItem(itemID string) error {
	qName := "delete_item"
	query, ok := d.queries[qName]
	if !ok {
		return fmt.Errorf("query missing %s", qName)
	}

	_, err := d.DB.Exec(query.Query, itemID)
	if err != nil {
		return fmt.Errorf("error deleting from items table: %w", err)
	}

	return nil
}

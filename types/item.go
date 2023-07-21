package types

type ShoppingListItem struct {
	ID             string `db:"id"`
	ShoppingListID string `db:"shopping_list_id"`
	ItemText       string `db:"item_text"`
	Purchased      bool   `db:"purchased"`
}

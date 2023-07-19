package types

import "time"

type ShoppingList struct {
	ID             string     `db:"id"`
	Title          string     `db:"title"`
	StoreLocation  string     `db:"store_location"`
	DueDate        *time.Time `db:"due_date"`
	TelegramChatID string     `db:"telegram_chat_id"`
}

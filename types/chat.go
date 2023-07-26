package types

type Chat struct {
	ID             string `db:"id" json:"id"`
	TelegramChatID int64  `db:"telegram_chat_id" json:"telegram_chat_id"`
}

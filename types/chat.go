package types

type Chat struct {
	ID             string `db:"id"`
	TelegramChatID int64  `db:"telegram_chat_id"`
}

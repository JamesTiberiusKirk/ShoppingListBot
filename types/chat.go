package types

type Chat struct {
	ID     string `db:"id"`
	ChatID int64  `db:"chat_id"`
}

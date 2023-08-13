package types

import (
	"time"

	"github.com/lib/pq"
)

type Journey struct {
	ID              string        `db:"id"`
	TelegramChatID  int64         `db:"telegram_chat_id"`
	ChatID          int64         `db:"chat_id"`
	Command         string        `db:"command"`
	Next            int           `db:"next"`
	RawContext      []byte        `db:"context"`
	MessagesCleanup pq.Int32Array `db:"messages_cleanup"`
	CreatedAt       time.Time     `db:"created_at"`
	UpdatedAt       time.Time     `db:"updated_at"`
}

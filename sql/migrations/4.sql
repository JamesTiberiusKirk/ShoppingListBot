ALTER TABLE usuals
ADD COLUMN chat_id INTEGER NOT NULL,
ADD FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE;

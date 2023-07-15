-- name: add_chat
INSERT INTO chats (chat_id) VALUES ($1);

-- name: get_chat
SELECT * FROM chats WHERE chat_id = $1;

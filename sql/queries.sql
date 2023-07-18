-- name: add_chat
INSERT INTO chats (telegram_chat_id) VALUES ($1);

-- name: get_chat
SELECT * FROM chats WHERE telegram_chat_id = $1;

-- name: add_list
INSERT INTO shopping_lists (chat_id, title, store_location, due_date)
VALUES ((SELECT id AS chat_id FROM chats WHERE telegram_chat_id = $1), $2,$3,$4);


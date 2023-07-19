-- name: add_chat
INSERT INTO chats (telegram_chat_id) VALUES ($1);

-- name: get_chat
SELECT * FROM chats WHERE telegram_chat_id = $1;

-- name: add_list
INSERT INTO shopping_lists (
    chat_id, title, store_location, due_date
)
VALUES (
    (
        SELECT id AS chat_id
        FROM chats WHERE telegram_chat_id = $1
    ), $2,$3,$4
);

-- name: get_lists
SELECT shopping_lists.id,  shopping_lists.title, shopping_lists.store_location, shopping_lists.due_date, chats.telegram_chat_id
FROM shopping_lists
JOIN chats ON shopping_lists.chat_id = chats.id
WHERE chats.telegram_chat_id = $1;

-- name: get_list_by_id
SELECT * 
FROM shopping_lists
WHERE id = $1 LIMIT 1;

-- name: add_items
INSERT INTO shopping_list_items (shopping_list_id, item_text, purchased)
VALUES (:list_id, :text, FALSE)

-- name: get_items_in_list
SELECT * FROM shopping_list_items WHERE shopping_list_id = $1

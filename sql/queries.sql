-- CHATS --

-- name: add_chat
INSERT INTO chats (telegram_chat_id) VALUES ($1);

-- name: get_chat
SELECT * FROM chats WHERE telegram_chat_id = $1;


-- LISTS --

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

-- name: delete_list_by_id
DELETE FROM shopping_lists
WHERE id = $1;




-- ITEMS --
--
-- name: add_items
INSERT INTO shopping_list_items (shopping_list_id, item_text, purchased)
VALUES (:list_id, :text, FALSE)

-- name: get_items_in_list
SELECT * FROM shopping_list_items WHERE shopping_list_id = $1 ORDER BY id;

-- name: get_unpurchased_items_in_list
SELECT * FROM shopping_list_items WHERE shopping_list_id = $1 AND purchased = FALSE ORDER BY id;

-- name: toggle_item_purchase
UPDATE shopping_list_items SET purchased = NOT purchased WHERE id = $1;

-- name: delete_item
DELETE FROM shopping_list_items WHERE id = $1;


-- JOURNIES --

-- name: upsert_journey_by_telegram_chat_id
INSERT INTO chat_journies (chat_id, command, next, context, messages_cleanup, created_at, updated_at)
VALUES (
    (
        SELECT id FROM chats WHERE telegram_chat_id = :telegram_chat_id
    ),
    :command,
    :next,
    :context,
    :messages_cleanup,
    (
        SELECT CURRENT_TIMESTAMP
    ),
    (
        SELECT CURRENT_TIMESTAMP
    )
)
ON CONFLICT (chat_id) DO UPDATE SET
        command = COALESCE(NULLIF(EXCLUDED.command, ''), chat_journies.command),
        next = COALESCE(NULLIF(EXCLUDED.next, -1), chat_journies.next),
        context = EXCLUDED.context,
        updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: get_chat_journey_by_telegram_chat_id
SELECT * FROM chat_journies 
JOIN chats ON chat_journies.chat_id = chats.id
WHERE chats.telegram_chat_id = $1

-- name: cleanup_chat_journies_by_telegram_chat_id
DELETE FROM chat_journies
WHERE chat_id = (
    SELECT id FROM chats WHERE telegram_chat_id = $1
);




-- USUALS --

-- name: get_all_usuals
SELECT * from usuals;

-- name: get_usual_by_id
SELECT * from usuals WHERE id = $1;

-- name: add_usuals;
INSERT INTO usuals (id, name, image_path, store) VALUES (
    :id,
    :name, 
    :image_path, 
    :store,
) ON CONFLICT (id) DO UPDATE SET 
    name = EXCLUDED.name,
    image_path = EXCLUDED.image_path,
    store = EXCLUDED.store
RETURNING *;

-- name: remove_usual
DELETE FROM usuals WHERE id = $1;

-- name: schema
CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    telegram_chat_id BIGINT UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS shopping_lists (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    store_location VARCHAR(255),
    due_date DATE,
    FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS shopping_list_items (
    id SERIAL PRIMARY KEY,
    shopping_list_id INTEGER NOT NULL,
    item_text VARCHAR(255) NOT NULL,
    purchased BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (shopping_list_id) REFERENCES shopping_lists (id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS chat_journies (
    id SERIAL PRIMARY KEY,
    chat_id BIGINT UNIQUE NOT NULL,
    command VARCHAR(255) NOT NULL,
    next INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    context JSONB,
    FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
);

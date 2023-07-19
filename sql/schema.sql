-- Create the tables

-- name: schema
CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    telegram_chat_id INTEGER UNIQUE NOT NULL
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

-- Example data

-- name: example_data
INSERT INTO chats (telegram_chat_id) VALUES (123456789);
INSERT INTO shopping_lists (chat_id, title, store_location, due_date)
VALUES (1, 'Groceries', 'Supermarket', '2023-07-31');
INSERT INTO shopping_list_items (shopping_list_id, item_text, purchased)
VALUES (1, 'Milk', FALSE),
       (1, 'Eggs', FALSE),
       (1, 'Bread', TRUE);

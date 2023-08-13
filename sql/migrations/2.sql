CREATE TABLE IF NOT EXISTS migrations (
    id SERIAL PRIMARY KEY,
    version INTEGER NOT NULL
);

INSERT INTO migrations (id, version) VALUES (1,2)
ON CONFLICT (id) DO UPDATE SET 
    version = COALESCE(NULLIF(EXCLUDED.version, 0), migrations.version);

ALTER TABLE chat_journies
    ADD COLUMN messages_cleanup INTEGER[];

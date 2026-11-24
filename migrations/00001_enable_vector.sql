-- +goose Up
CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;

-- +goose Down
-- The extension is shared by schemas. Keep it when rolling back application tables.
SELECT 1;

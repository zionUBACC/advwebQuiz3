-- Filename: migrations/000003_add_todo_indexes.up.sql
CREATE INDEX IF NOT EXISTS todo_item_idx ON todolist USING GIN(to_tsvector('simple', item));
CREATE INDEX IF NOT EXISTS todo_description_idx ON todolist USING GIN(to_tsvector('simple', description));

-- Filename: migrations/000003_add_todo_indexes.down.sql
DROP INDEX If EXISTS todo_item_idx;
DROP INDEX If EXISTS todo_description_idx;

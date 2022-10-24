CREATE TABLE IF NOT EXISTS todolist (
    id serial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    item TEXT NOT NULL,
    description TEXT
);
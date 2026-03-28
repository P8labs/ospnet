-- +goose Up
CREATE TABLE containers (
    id TEXT PRIMARY KEY,
    image TEXT NOT NULL,
    name TEXT NOT NULL,
    port INTEGER NOT NULL,
    status TEXT NOT NULL,
    docker_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE containers;
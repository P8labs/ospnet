-- +goose Up
CREATE TABLE nodes (
    id TEXT PRIMARY KEY,

    name TEXT NOT NULL,
    hostname TEXT NOT NULL,
    ip TEXT NOT NULL,
    cpu INTEGER NOT NULL,
    memory INTEGER NOT NULL,
    arch TEXT NOT NULL,
    region TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    last_seen TIMESTAMP NOT NULL,
    labels TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE onboarding_tokens (
    token TEXT PRIMARY KEY,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tokens_expires ON onboarding_tokens(expires_at);

CREATE INDEX idx_nodes_last_seen ON nodes(last_seen);
CREATE INDEX idx_nodes_region ON nodes(region);
CREATE INDEX idx_nodes_type ON nodes(type);

-- +goose Down
DROP TABLE nodes;
DROP TABLE onboarding_tokens;
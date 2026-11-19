-- +goose Up
CREATE TABLE settings (
    id INTEGER PRIMARY KEY,
    key TEXT NOT NULL,
    value TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE UNIQUE INDEX idx_settings_key ON settings(key);

INSERT INTO settings (key, value, created_at, updated_at)
VALUES ('name', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO settings (key, value, created_at, updated_at)
VALUES ('count', '0', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- +goose Down
DROP TABLE settings;

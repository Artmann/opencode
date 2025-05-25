-- +goose Up
-- +goose StatementBegin
-- Message History
CREATE TABLE IF NOT EXISTS message_history (
    id TEXT PRIMARY KEY,
    session_id TEXT REFERENCES sessions(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f000Z', 'now')),
    UNIQUE(session_id, message)
);

CREATE INDEX message_history_session_id_idx ON message_history(session_id);
CREATE INDEX message_history_created_at_idx ON message_history(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS message_history_created_at_idx;
DROP INDEX IF EXISTS message_history_session_id_idx;
DROP TABLE IF EXISTS message_history;
-- +goose StatementEnd
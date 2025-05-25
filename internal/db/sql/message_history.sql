-- Message History queries

-- name: CreateMessageHistory :one
INSERT INTO message_history (id, session_id, message)
VALUES (?, ?, ?)
ON CONFLICT(session_id, message) DO NOTHING
RETURNING *;

-- name: GetMessageHistoryBySession :many
SELECT * FROM message_history
WHERE session_id = ?
ORDER BY created_at ASC;

-- name: DeleteMessageHistoryBySession :exec
DELETE FROM message_history
WHERE session_id = ?;

-- name: DeleteMessageHistory :exec
DELETE FROM message_history
WHERE id = ?;
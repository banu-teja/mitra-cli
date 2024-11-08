-- name: CreateCommandRequest :one
INSERT INTO command_requests (
    input_content,
    command_type
) VALUES (
    ?, ?
) RETURNING *;

-- name: GetCommandRequest :one
SELECT * FROM command_requests
WHERE id = ? LIMIT 1;

-- name: ListCommandRequests :many
SELECT * FROM command_requests
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateCommandRequest :one
UPDATE command_requests
SET input_content = ?,
    command_type = ?
WHERE id = ?
RETURNING *;

-- name: DeleteCommandRequest :exec
DELETE FROM command_requests
WHERE id = ?;

-- name: GetCommandRequestsByType :many
SELECT * FROM command_requests
WHERE command_type = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

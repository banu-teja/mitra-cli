-- name: CreateSubCommand :one
INSERT INTO sub_commands (
    request_id,
    command,
    command_output,
    command_status,
    execution_order
) VALUES (
    ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetSubCommand :one
SELECT * FROM sub_commands
WHERE id = ? LIMIT 1;

-- name: ListSubCommands :many
SELECT * FROM sub_commands
WHERE request_id = ?
ORDER BY execution_order ASC;

-- name: UpdateSubCommandOutput :one
UPDATE sub_commands
SET command_output = ?
WHERE id = ?
RETURNING *;


-- name: DeleteSubCommand :exec
DELETE FROM sub_commands
WHERE id = ?;

-- name: GetSubCommandsByStatus :many
SELECT * FROM sub_commands
WHERE command_status = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateSubCommandStatus :one
UPDATE sub_commands
SET command_status = ?,
    command_output = ?
WHERE id = ?
RETURNING *;

-- name: GetSubCommandsWithRequest :many
SELECT sc.*, cr.input_content, cr.command_type
FROM sub_commands sc
JOIN command_requests cr ON sc.request_id = cr.id
WHERE sc.request_id = ?
ORDER BY sc.execution_order ASC;

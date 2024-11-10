-- name: GetLastNEntries :many
WITH last_n_requests AS (
    SELECT id, input_content, command_type, created_at
    FROM command_requests
    ORDER BY created_at DESC
    LIMIT ?
)
SELECT 
    r.id AS request_id,
    r.input_content,
    r.command_type,
    r.created_at AS request_created_at,
    s.id AS subcommand_id,
    s.command,
    s.command_output,
    s.command_status,
    s.execution_order,
    s.created_at AS subcommand_created_at
FROM last_n_requests r
LEFT JOIN sub_commands s ON r.id = s.request_id
ORDER BY r.created_at DESC, s.execution_order;

-- schema.sql

CREATE TABLE command_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    input_content TEXT NOT NULL,
    command_type TEXT NOT NULL CHECK(command_type IN ('revert', 'history', 'question')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sub_commands (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id INTEGER NOT NULL,
    command TEXT NOT NULL,
    command_output TEXT NOT NULL,
    command_status TEXT NOT NULL CHECK(command_status IN ('success', 'failure', 'in_progress')),
    execution_order INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_id) REFERENCES command_requests(id)
);

-- Index for faster lookups
CREATE INDEX idx_sub_commands_request_id ON sub_commands(request_id);

-- +migrate Up
CREATE TABLE IF NOT EXISTS boards (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_boards_project ON boards(project_id);
CREATE INDEX idx_boards_position ON boards(project_id, position);

-- +migrate Down
DROP TABLE IF EXISTS boards;

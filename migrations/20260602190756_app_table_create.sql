-- +goose Up
CREATE TABLE IF NOT EXISTS app (
    id VARCHAR(255) PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    repo_url TEXT NOT NULL,
    branch VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    env_vars JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    uodated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS app;

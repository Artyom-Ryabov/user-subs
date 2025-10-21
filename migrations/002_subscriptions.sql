-- +goose Up
ALTER TABLE subscriptions ADD COLUMN ended_at TIMESTAMP;

-- +goose Down
ALTER TABLE subscriptions DROP COLUMN ended_at;

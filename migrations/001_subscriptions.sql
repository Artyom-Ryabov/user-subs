-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INT NOT NULL DEFAULT 0,
    user_id UUID NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO subscriptions (
    service_name,
    price,
    user_id
) VALUES (
    '[SERVICE NAME 1]',
    420,
    '60601fee-2bf1-4721-ae6f-7636e79a0cba'
), (
    '[SERVICE NAME 2]',
    69,
    '60601fee-2bf1-4721-ae6f-7636e79a0cbb'
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE subscriptions;

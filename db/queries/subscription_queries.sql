-- name: GetSubs :many
SELECT * FROM subscriptions;

-- name: GetSub :one
SELECT * FROM subscriptions WHERE id = $1;

-- name: GetUserSubs :many
SELECT * FROM subscriptions WHERE user_id = $1;

-- name: AddSub :one
INSERT INTO subscriptions (
    service_name,
    price,
    user_id,
    started_at
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING id;

-- name: UpdateSub :one
UPDATE subscriptions SET
    service_name = $1,
    price = $2,
    user_id = $3,
    started_at = $4,
    updated_at = $5
WHERE id = $6 RETURNING id;

-- name: DeleteSub :one
DELETE FROM subscriptions WHERE id = $1 RETURNING id;

-- name: DeleteUserSubs :many
DELETE FROM subscriptions WHERE user_id = $1 RETURNING id;

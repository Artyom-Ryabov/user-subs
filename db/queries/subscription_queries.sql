-- name: GetSubs :many
SELECT
    *
FROM
    subscriptions;
-- name: AddSub :one
INSERT INTO subscriptions(
    service_name,
    price,
    user_id,
    started_at
)VALUES(
    $1,
    $2,
    $3,
    $4
) RETURNING id;

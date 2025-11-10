-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
VALUES (sqlc.arg(user_id)::uuid, sqlc.arg(service_name)::TEXT, sqlc.arg(price)::INT, sqlc.arg(start_date)::DATE,
        sqlc.narg(end_date)::DATE)
RETURNING id, user_id, service_name, price, start_date, end_date;

-- name: GetSubscriptionById :one
SELECT id, user_id, service_name, price, start_date, end_date
FROM subscriptions
WHERE id = sqlc.arg(subscription_id);

-- name: AllSubscriptions :many
SELECT id, user_id, service_name, price, start_date, end_date
FROM subscriptions;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET service_name = COALESCE(sqlc.narg(service_name)::TEXT, service_name),
    price        = COALESCE(sqlc.narg(price)::INT, price),
    end_date     = COALESCE(sqlc.narg(end_date)::DATE, end_date)
WHERE id = sqlc.arg(subscription_id)
RETURNING id, user_id, service_name, price, start_date, end_date;

-- name: DeleteSubscription :exec
DELETE
FROM subscriptions
WHERE id = sqlc.arg(subscription_id);

-- name: GetSumSubscriptions :one
SELECT COALESCE(SUM(price), 0)::INT
FROM subscriptions
WHERE start_date <= sqlc.arg(end_date)::DATE
  AND (end_date IS NULL OR end_date >= sqlc.arg(start_date)::DATE)
  AND (sqlc.narg(service_name)::TEXT IS NULL OR service_name ILIKE sqlc.narg(service_name)::TEXT)
  AND (sqlc.narg(user_id)::uuid IS NULL OR user_id = sqlc.narg(user_id)::uuid);

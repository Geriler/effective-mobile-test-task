-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions
(
    id           uuid DEFAULT gen_random_uuid() NOT NULL
        CONSTRAINT subscriptions_pk PRIMARY KEY,
    user_id      uuid                           NOT NULL,
    service_name TEXT                           NOT NULL,
    price        INTEGER                        NOT NULL,
    start_date   DATE                           NOT NULL,
    end_date     DATE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd

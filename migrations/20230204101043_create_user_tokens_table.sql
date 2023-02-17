-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_tokens
(
    id                       serial primary key not null,
    pair_uuid                varchar(50) unique not null,
    user_id                  bigint             not null,
    access_token             varchar(255) unique,
    refresh_token            varchar(255) unique,
    is_valid                 bool               NOT NULL DEFAULT true,
    access_token_expired_at  timestamp(0)       NOT NULL,
    refresh_token_expired_at timestamp(0)       NOT NULL,
    updated_at               timestamp(0)       NOT NULL DEFAULT now(),
    created_at               timestamp(0)       NOT NULL DEFAULT now(),

    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user_tokens";
-- +goose StatementEnd

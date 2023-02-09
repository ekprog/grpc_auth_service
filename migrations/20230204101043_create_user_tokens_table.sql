-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_tokens
(
    id         serial primary key  not null,
    user_id    bigint              not null,
    token      varchar(255) unique not null,
    is_valid   bool                NOT NULL DEFAULT true,
    updated_at timestamp(0)        NOT NULL DEFAULT now(),
    created_at timestamp(0)        NOT NULL DEFAULT now(),
    expired_at timestamp(0)        NOT NULL,

    constraint fk_user_id foreign key (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user_tokens";
-- +goose StatementEnd

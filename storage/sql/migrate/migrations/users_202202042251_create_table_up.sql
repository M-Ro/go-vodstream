CREATE TABLE users (
    id          BIGSERIAL   PRIMARY KEY,
    username    TEXT        NOT NULL,
    email       TEXT        NOT NULL,
    password    TEXT,
    publish_key TEXT,
    can_publish BOOL,
    can_stream  BOOL,
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP
);
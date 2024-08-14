CREATE TABLE IF NOT EXISTS users
(
    user_id  uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    username varchar(16)      NOT NULL,
    password varchar(32)      NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens
(
    user_id       uuid PRIMARY KEY NOT NULL REFERENCES users (user_id),
    refresh_token text             NOT NULL
);
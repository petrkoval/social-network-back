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

CREATE TABLE IF NOT EXISTS channels
(
    channel_id  uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id     uuid             NOT NULL REFERENCES users (user_id),
    created_at  timestamp        NOT NULL DEFAULT now(),
    title       varchar(32)      NOT NULL,
    description varchar(256)              DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    post_id    uuid PRIMARY KEY NOT NULL             DEFAULT uuid_generate_v4(),
    user_id    uuid REFERENCES users (user_id)       DEFAULT NULL,
    channel_id uuid REFERENCES channels (channel_id) DEFAULT NULL,
    created_at timestamp        NOT NULL             DEFAULT now(),
    content    text                                  DEFAULT NULL,
    images     text[]                                DEFAULT NULL
);
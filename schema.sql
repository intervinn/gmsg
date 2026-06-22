
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE guilds (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    owner_id BIGINT REFERENCES users(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE channels (
    id BIGINT PRIMARY KEY,
    name TEXT,
    guild_id BIGINT REFERENCES guilds(id) ON DELETE CASCADE,
    type SMALLINT NOT NULL
);

CREATE TABLE messages (
    id BIGINT PRIMARY KEY,
    content TEXT NOT NULL,
    author_id BIGINT REFERENCES users(id) NOT NULL,
    guild_id BIGINT REFERENCES guilds(id),
    channel_id BIGINT REFERENCES channels(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE guild_members (
    guild_id BIGINT REFERENCES guilds(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (guild_id, user_id)
);
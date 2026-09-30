CREATE TABLE IF NOT EXISTS chats (
    id BIGINT PRIMARY KEY,
    play_type INT NOT NULL DEFAULT 0,
    admin_play BOOLEAN NOT NULL DEFAULT FALSE,
    admin_mode VARCHAR(50) NOT NULL DEFAULT 'everyone',
    cmd_delete BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS playlists (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_playlists_user_id ON playlists(user_id);

CREATE TABLE IF NOT EXISTS playlist_songs (
    playlist_id VARCHAR(50) NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    track_id VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    name VARCHAR(255) NOT NULL,
    duration INT NOT NULL,
    platform VARCHAR(50) NOT NULL,
    PRIMARY KEY (playlist_id, track_id)
);

CREATE TABLE IF NOT EXISTS assistants (
    chat_id BIGINT PRIMARY KEY,
    num INT NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_users (
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS blacklisted_chats (
    chat_id BIGINT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS blacklisted_users (
    user_id BIGINT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS chat_languages (
    chat_id BIGINT PRIMARY KEY,
    lang_code VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS bot_settings (
    key VARCHAR(50) PRIMARY KEY,
    value TEXT NOT NULL
);

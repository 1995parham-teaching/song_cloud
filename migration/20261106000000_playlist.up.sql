CREATE TABLE IF NOT EXISTS playlist
(
    id         SERIAL PRIMARY KEY,
    name       varchar(255) NOT NULL CHECK ( LENGTH(TRIM(name)) > 0 ),
    owner      varchar(255) NOT NULL,
    created_at timestamp    NOT NULL DEFAULT now(),
    CONSTRAINT fk_playlist_owner FOREIGN KEY (owner) REFERENCES users (username) ON DELETE CASCADE,
    CONSTRAINT uq_playlist_owner_name UNIQUE (owner, name)
);

-- songs of a playlist, ordered by the time they were added.
CREATE TABLE IF NOT EXISTS playlist_song
(
    playlist_id int       NOT NULL,
    song_id     int       NOT NULL,
    added_at    timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (playlist_id, song_id),
    CONSTRAINT fk_playlist_song_playlist FOREIGN KEY (playlist_id) REFERENCES playlist (id) ON DELETE CASCADE,
    CONSTRAINT fk_playlist_song_song FOREIGN KEY (song_id) REFERENCES song (id) ON DELETE CASCADE
);

-- users that a playlist is shared with (read-only access).
CREATE TABLE IF NOT EXISTS playlist_share
(
    playlist_id int          NOT NULL,
    username    varchar(255) NOT NULL,
    shared_at   timestamp    NOT NULL DEFAULT now(),
    PRIMARY KEY (playlist_id, username),
    CONSTRAINT fk_playlist_share_playlist FOREIGN KEY (playlist_id) REFERENCES playlist (id) ON DELETE CASCADE,
    CONSTRAINT fk_playlist_share_username FOREIGN KEY (username) REFERENCES users (username) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_playlist_share_username ON playlist_share (username);

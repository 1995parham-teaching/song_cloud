CREATE TABLE IF NOT EXISTS "like"
(
    id              SERIAL primary key,
    username    varchar(255),
    song_id     INT,
    CONSTRAINT FK_song FOREIGN KEY (song_id) REFERENCES song(id),
    CONSTRAINT FK_username FOREIGN KEY (username) REFERENCES users(username)
);

CREATE OR REPLACE PROCEDURE like_song(username_in varchar(255), id_in integer) AS
$$
  BEGIN
    INSERT INTO "like" (username, song_id) VALUES (username_in, id_in);
  END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION like_log() RETURNS trigger as
$$
  BEGIN
    INSERT INTO "log"(status_code, log_message, "time")
        VALUES (0, CONCAT(new.username, ' liked this song: ', new.song_id), now());
    RETURN NEW;
  END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER user_like AFTER INSERT on "like" FOR EACH ROW EXECUTE PROCEDURE like_log();

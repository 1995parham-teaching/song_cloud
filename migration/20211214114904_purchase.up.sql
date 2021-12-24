CREATE TABLE IF NOT EXISTS purchase
(
    id SERIAL primary key,
    username   varchar(255),
    song_id     int,
    purchased_date timestamp,
    purchased_price int,
    CONSTRAINT FK_song FOREIGN KEY (song_id) REFERENCES song(id),
    CONSTRAINT FK_username FOREIGN KEY (username) REFERENCES users(username)
);

CREATE OR REPLACE FUNCTION isPaid(username_in varchar(255), song_id_in int) RETURNS BOOLEAN as
$$
  BEGIN
    IF EXISTS(SELECT FROM purchase WHERE purchase.username = username_in and purchase.song_id = song_id_in) THEN
      RETURN TRUE;
    END IF;
    RETURN FALSE;
END;
$$
LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION purchase_log() RETURNS trigger as
$$
DECLARE
    song_price integer;
BEGIN
  SELECT price INTO song_price FROM song WHERE song.id = new.song_id;
  INSERT INTO "log"(status_code, log_message, "time")
      VALUES (0, CONCAT(new.username, ' purchased this song ', new.song_id, ' by paying ', song_price), now());
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER tr_purchase_log AFTER INSERT on "purchase" FOR EACH ROW EXECUTE PROCEDURE purchase_log();

create role bank login password 'bank';
grant select on purchase to bank;

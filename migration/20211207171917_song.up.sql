CREATE TABLE IF NOT EXISTS song
(
    id              SERIAL primary key,
    file            varchar(255),
    name            varchar(255),
    production_year int,
    explanation     varchar(255),
    view            int default 0,
    price           int DEFAULT 0,
    score           numeric default 0
);

CREATE OR REPLACE FUNCTION free(id integer) RETURNS boolean AS
$$
    BEGIN
      IF exists(select from song where song.id = id and price = 0) THEN
        RETURN TRUE;
      END IF;
      RETURN FALSE;
    END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE increase_view(id integer) AS
$$
  BEGIN
    UPDATE song SET view = view + 1 WHERE song.id = id;
  END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE play(id integer, username varchar(255)) AS
$$
  DECLARE
    song_price integer;
  BEGIN
    SELECT price INTO song_price FROM song WHERE song.id = id;
    IF free(id) THEN
      CALL increase_view(id);
    END IF;
    IF premium_user_validation(username) THEN
      CALL increase_view(id);
    END IF;
    IF pay(username, song_price) THEN
      CALL increase_view(id);
    END IF;
  END;
$$
LANGUAGE plpgsql;

-- CREATE TRIGGER watch_create BEFORE INSERT on watch FOR EACH ROW EXECUTE PROCEDURE pay_for_watch();

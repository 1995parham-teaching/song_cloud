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

CREATE OR REPLACE FUNCTION free(id_in integer) RETURNS boolean AS
$$
    BEGIN
      IF exists(select from song where song.id = id_in and price = 0) THEN
        RETURN TRUE;
      END IF;
      RETURN FALSE;
    END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE increase_view(id_in integer) AS
$$
  BEGIN
    UPDATE song SET view = view + 1 WHERE song.id = id_in;
  END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE play(id_in integer, username_in varchar(255)) AS
$$
  BEGIN
    IF free(id_in) THEN
      CALL increase_view(id_in);
      RETURN;
    END IF;
    IF premium_user_validation(username_in) THEN
      CALL increase_view(id_in);
      RETURN;
    END IF;
    IF isPaid(username_in, id_in) THEN
      CALL increase_view(id_in);
    END IF;
    RAISE EXCEPTION 'user % cannot play %', username_in, id_in;
  END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_song(
    file_in            varchar(255),
    name_in            varchar(255),
    production_year_in int,
    explanation_in     varchar(255) DEFAULT '',
    price_in           int DEFAULT 0,
    score_in           numeric DEFAULT 0) RETURNS integer as
$$
  BEGIN
    INSERT INTO public.song(
      file, name, production_year, explanation, price, score)
    VALUES (file_in, name_in, production_year_in, explanation_in, price_in, score_in);
    RETURN currval('song_id_seq');
  END;
$$
LANGUAGE plpgsql;


-- CREATE TRIGGER watch_create BEFORE INSERT on watch FOR EACH ROW EXECUTE PROCEDURE pay_for_watch();

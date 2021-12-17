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
    END IF;
    IF premium_user_validation(username_in) THEN
      CALL increase_view(id_in);
    END IF;
    IF isPaid(username_in, id_in) THEN  
      CALL increase_view(id_in);
    END IF;
  END;
$$
LANGUAGE plpgsql;

-- CREATE TRIGGER watch_create BEFORE INSERT on watch FOR EACH ROW EXECUTE PROCEDURE pay_for_watch();

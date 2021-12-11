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

CREATE OR REPLACE free()
$$
    BEGIN
      IF (select count(1) from song where NEW.id = id) THEN
        RETURN TRUE
      END IF;
      RETURN FALSE;
    END;
$$
LANGUAGE 'plpgsql'

CREATE OR REPLACE increase_view()
$$
  BEGIN
    UPDATE song SET view = view + 1 WHERE id = NEW.id
    RETURN NULL;
  END;
$$
LANGUAGE 'plpgsql'


CREATE OR REPLACE FUNCTION play() RETURNS trigger as
$$
  BEGIN
    IF free() THEN
      increase_view()
      RETURN TRUE;
    END IF;
    IF premium_user_validation() THEN
      increase_view()
      RETURN TRUE;
    END IF;
    IF pay() THEN
      increase_view()
      RETURN TRUE;
    END IF;
    RETURN NULL;
  END;
$$
LANGUAGE 'plpgsql';

CREATE TRIGGER watch_create BEFORE INSERT on watch FOR EACH ROW EXECUTE PROCEDURE pay_for_watch();

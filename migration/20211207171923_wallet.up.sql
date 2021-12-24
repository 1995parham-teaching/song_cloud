CREATE TABLE IF NOT EXISTS wallet
(
    id SERIAL primary key ,
    username   varchar(255),
    credit     int default 0,
);

CREATE OR REPLACE PROCEDURE pay(username_in varchar(255), song_id_in int) as
$$
DECLARE
    song_price integer;
BEGIN
    SELECT price INTO song_price FROM song WHERE song.id = song_id_in;
    IF song_price = 0 THEN
      RETURN;
    END IF;
    IF (premium_user_validation(username_in)) THEN
      RETURN;
    END IF;
    IF (isPaid(username_in, song_id_in)) THEN
      RETURN;
    END IF;
    IF (free(song_id_in)) THEN
      RETURN;
    END IF;
    IF EXISTS (SELECT FROM wallet WHERE wallet.username = username_in and wallet.credit >= song_price) THEN
      UPDATE wallet SET credit = credit - song_price WHERE wallet.username = username_in;
      INSERT INTO purchase(username, song_id) VALUES (username_in, song_id_in);
      COMMIT;
    ELSE
      RAISE EXCEPTION 'not enough credit (%) for user % to buy %', song_price, username_in, song_id_in;
    END IF;
END;
$$
LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION charge_credit(username_in varchar(255), amount_in INTEGER) RETURNS boolean as
$$
BEGIN
    IF EXISTS(SELECT FROM users WHERE users.username = username_in) THEN
      UPDATE wallet SET credit = credit + amount_in WHERE wallet.username = username_in;
      RETURN true;
    END IF;
    RETURN false;
END;
$$
LANGUAGE plpgsql;

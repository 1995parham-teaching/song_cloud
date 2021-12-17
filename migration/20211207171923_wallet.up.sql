CREATE TABLE IF NOT EXISTS wallet
(
    id SERIAL primary key ,
    username   varchar(255),
    credit     int default 0,
    explanation varchar(255)
);

CREATE OR REPLACE FUNCTION pay(username_in varchar(255), song_id_in int) RETURNS BOOLEAN as
$$
DECLARE
    song_price integer;
BEGIN
    SELECT price INTO song_price FROM song WHERE song.id = song_id_in;
    IF song_price = 0 THEN
      RETURN FALSE;
    END IF;
    IF (NOT isPaid(username_in, song_id_in)) AND EXISTS (SELECT FROM wallet WHERE wallet.username = username_in and wallet.credit >= song_price) THEN
      UPDATE wallet SET credit = credit - song_price WHERE wallet.username = username_in;
      INSERT INTO purchase(username, song_id, paid) VALUES (username_in, song_id_in, song_price);
      RETURN TRUE;
    END IF;
    RETURN FALSE;
END;
$$
LANGUAGE plpgsql;

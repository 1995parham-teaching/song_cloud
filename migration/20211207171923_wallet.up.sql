CREATE TABLE IF NOT EXISTS wallet
(
    id int primary key ,
    username   varchar(255),
    credit     int default 0,
    explanation varchar(255)
);

CREATE OR REPLACE pay(price float) RETURNS BOOLEAN as
$$
BEGIN
    IF (SELECT COUNT(1) FROM wallet WHERE username=NEW.username and credit >= NEW.credit)
      UPDATE wallet SET credit = credit - price WHERE username = NEW.username;
      RETURN TRUE;
    END IF;
    RETURN FALSE;
END;
$$
LANGUAGE 'plpgsql'
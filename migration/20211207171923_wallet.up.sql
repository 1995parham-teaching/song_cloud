CREATE TABLE IF NOT EXISTS wallet
(
    id int primary key ,
    username   varchar(255),
    credit     int default 0,
    explanation varchar(255)
);

CREATE OR REPLACE FUNCTION pay(username varchar(255), price int) RETURNS BOOLEAN as
$$
BEGIN
    IF EXISTS(SELECT FROM wallet WHERE wallet.username = username and wallet.credit >= price)
      UPDATE wallet SET credit = credit - price WHERE wallet.username = username;
      RETURN TRUE;
    END IF;
    RETURN FALSE;
END;
$$
LANGUAGE plpgsql

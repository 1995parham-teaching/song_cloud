CREATE TABLE IF NOT EXISTS users
(
    username     varchar(255) primary key,
    password     varchar(255) check ( LENGTH(password) >= 8 and (password similar to '%[0-9A-Za-z]%') ),
    first_name   varchar(255),
    last_name    varchar(255),
    email        varchar(255),
    premium_till date,
    score        int default 0
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_email ON users (email);

CREATE OR REPLACE FUNCTION premium_user_validation() RETURNS trigger as
$$
BEGIN
    IF (select count(1) from users where username = NEW.username and premium_till > now()) = 1 THEN
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$
    LANGUAGE 'plpgsql';
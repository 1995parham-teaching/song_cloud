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

CREATE OR REPLACE FUNCTION premium_user_validation(username varchar(255)) RETURNS boolean as
$$
BEGIN
    IF EXISTS(SELECT FROM users WHERE users.username = username AND users.premium_till > now()) THEN
      RETURN true;
    END IF;
    RETURN false;
END;
$$
LANGUAGE plpgsql;

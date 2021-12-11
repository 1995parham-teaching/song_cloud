CREATE TABLE IF NOT EXISTS users
(
    username        varchar(255) primary key,
    password        varchar(255) check ( LENGTH(password) >= 8 and (password similar to '%[0-9A-Za-z]%') ),
    first_name      varchar(255),
    last_name       varchar(255),
    email           varchar(255),
    phone           varchar(13),
    national_number varchar(10),
    address         varchar(255),
    score           int default 0
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_email ON users(email);

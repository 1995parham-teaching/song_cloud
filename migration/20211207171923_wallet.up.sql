CREATE TABLE IF NOT EXISTS wallet
(
    id int primary key ,
    username   varchar(255),
    credit     int default 0,
    explanation varchar(255)
)
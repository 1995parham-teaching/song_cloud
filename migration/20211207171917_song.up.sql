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
)
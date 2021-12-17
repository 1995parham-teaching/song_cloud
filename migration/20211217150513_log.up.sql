CREATE TABLE IF NOT EXISTS log
(
    id              SERIAL primary key,
    status_code INTEGER,
    log_message varchar(500),
    "time"  timestamp  
);


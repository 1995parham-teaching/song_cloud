CREATE TABLE IF NOT EXISTS log
(
  id              SERIAL primary key,
  status_code INTEGER,
  log_message varchar(500),
  "time"  timestamp
);

create role logman login password 'logman';
grant select on log to logman;

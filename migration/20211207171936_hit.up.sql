CREATE TABLE IF NOT EXISTS hits
(
    song  int,
    views int,
    likes int,
    FOREIGN KEY (song) REFERENCES song (id)
);
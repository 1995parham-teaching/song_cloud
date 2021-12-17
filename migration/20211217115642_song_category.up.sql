CREATE TABLE IF NOT EXISTS song_category
(
    id              SERIAL primary key,
    song_id     INT,
    category_id INT,
    CONSTRAINT FK_song FOREIGN KEY (song_id) REFERENCES song(id),
    CONSTRAINT FK_category FOREIGN KEY (category_id) REFERENCES category(id)
);
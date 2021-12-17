CREATE TABLE IF NOT EXISTS song_category
(
    id              SERIAL primary key,
    song_id     INT,
    category_id INT,
    CONSTRAINT FK_song FOREIGN KEY (song_id) REFERENCES song(id),
    CONSTRAINT FK_category FOREIGN KEY (category_id) REFERENCES category(id)
);


CREATE OR REPLACE FUNCTION get_song_category(song_id_in integer) 
RETURNS TABLE(
    category_id INT,
    category_name varchar(255)
)
AS
$$
  BEGIN
    RETURN QUERY
        (SELECT b.id category_id, b.category_name 
        FROM song_category a, category b
        WHERE song_id_in = a.song_id
            AND a.category_id = b.id);
  END;
$$
LANGUAGE plpgsql;
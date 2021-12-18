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


CREATE OR REPLACE PROCEDURE assign_category_to_song(
    song_id_in  int,
    category_id_in int) as
$$
  BEGIN
    IF (NOT EXISTS (SELECT FROM public.song_category WHERE song_id = song_id_in and category_id = category_id_in)) THEN
      INSERT INTO public.song_category(song_id, category_id)
      VALUES (song_id_in, category_id_in);
    END IF;
  END;
$$
LANGUAGE plpgsql;
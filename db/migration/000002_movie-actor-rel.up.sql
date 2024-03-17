CREATE TABLE movie_actors (
    movie_id INT REFERENCES movies(id),
    actor_id INT REFERENCES actors(id),
    PRIMARY KEY (movie_id, actor_id)
);

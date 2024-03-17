INSERT INTO movies (name, description, release_date, rating)
VALUES
    ('Inception', 'A mind-bending thriller directed by Christopher Nolan.', '2010-07-16', 8.8),
    ('The Shawshank Redemption', 'A classic drama film based on the novella by Stephen King.', '1994-09-23', 9.3),
    ('The Godfather', 'An epic crime film directed by Francis Ford Coppola.', '1972-03-24', 9.2),
    ('The Dark Knight', 'A superhero film directed by Christopher Nolan.', '2008-07-18', 9.0),
    ('Pulp Fiction', 'A black comedy crime film directed by Quentin Tarantino.', '1994-05-21', 8.9),
    ('Forrest Gump', 'A comedy-drama film directed by Robert Zemeckis.', '1994-07-06', 8.8),
    ('The Matrix', 'A science fiction action film directed by the Wachowskis.', '1999-03-31', 8.7);
INSERT INTO actors (name, gender, birthday)
VALUES
    ('Leonardo DiCaprio', 'male', '1974-11-11'),
    ('Morgan Freeman', 'male', '1937-06-01'),
    ('Marlon Brando', 'male', '1924-04-03'),
    ('Tom Hanks', 'male', '1956-07-09'),
    ('Keanu Reeves', 'male', '1964-09-02'),
    ('Samuel L. Jackson', 'male', '1948-12-21'),
    ('Uma Thurman', 'female', '1970-04-29'),
    ('Robert De Niro', 'male', '1943-08-17'),
    ('Al Pacino', 'male', '1940-04-25'),
    ('Denzel Washington', 'male', '1954-12-28'),
    ('Meryl Streep', 'female', '1949-06-22'),
    ('Brad Pitt', 'male', '1963-12-18'),
    ('Angelina Jolie', 'female', '1975-06-04'),
    ('Johnny Depp', 'male', '1963-06-09'),
    ('Nicole Kidman', 'female', '1967-06-20');

INSERT INTO users (username, hashed_password, password_changed_at, role, created_at)
VALUES
    ('admin', 'qwerty', NOW(), 'administrator', NOW());


INSERT INTO movie_actors (movie_id, actor_id)
VALUES
    (1, 1), -- Inception - Leonardo DiCaprio
    (1, 2), -- Inception - Morgan Freeman
    (2, 1), -- The Shawshank Redemption - Leonardo DiCaprio
    (2, 3), -- The Shawshank Redemption - Marlon Brando
    (3, 1), -- The Godfather - Leonardo DiCaprio
    (3, 3), -- The Godfather - Marlon Brando
    (4, 4), -- The Dark Knight - Tom Hanks
    (4, 5), -- The Dark Knight - Keanu Reeves
    (5, 6), -- Pulp Fiction - Samuel L. Jackson
    (5, 7); -- Pulp Fiction - Uma Thurman


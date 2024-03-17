-- name: CreateMovie :one
INSERT INTO movies (
  name,
  description,
  release_date,
  rating
) VALUES 
  ($1, $2, $3, $4) RETURNING *;

-- name: UpdateMovie :one
UPDATE movies
SET name = $2,
  description = $3,
  rating = $4,
  release_date = $5
WHERE id = $1
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;

-- name: GetMoviesSortedByRating :many
SELECT *
FROM movies
ORDER BY rating DESC;

-- name: GetMoviesSortedByName :many
SELECT *
FROM movies
ORDER BY name;

-- name: GetMoviesByReleaseDate :many
SELECT *
FROM movies
ORDER BY release_date DESC;

-- name: GetMoviesByNameFragment :many
SELECT *
FROM movies
WHERE name LIKE '%' || $1 || '%';

-- name: GetMoviesByActorFragment :many
SELECT m.*
FROM movies m
JOIN movie_actors ma ON m.id = ma.movie_id
JOIN actors a ON ma.actor_id = a.id
WHERE a.name LIKE '%' || $1 || '%';

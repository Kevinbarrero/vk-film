-- name: CreateActor :one
INSERT INTO actors (
  name,
  gender,
  birthday
) VALUES 
  ($1, $2, $3) RETURNING *;

-- name: UpdateActor :one
UPDATE actors
SET name = $2,
  gender = $3,
  birthday = $4
WHERE id = $1
RETURNING *;

-- name: DeleteActor :exec
DELETE FROM actors
WHERE id = $1;

-- name: GetActorMoviesList :many
SELECT a.name AS actor_name, m.name AS movie_name
FROM actors a
JOIN movie_actors ma ON a.id = ma.actor_id
JOIN movies m ON ma.movie_id = m.id
ORDER BY a.name, m.name;

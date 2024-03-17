-- name: CreateActor :one
INSERT INTO actors (
  name,
  gender,
  birthday
) VALUES 
  ($1, $2, $3) RETURNING *;


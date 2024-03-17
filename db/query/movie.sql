-- name: CreateMovie :one
INSERT INTO movies (
  name,
  description,
  release_date,
  rating
) VALUES 
  ($1, $2, $3, $4) RETURNING *;





// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: actor.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createActor = `-- name: CreateActor :one
INSERT INTO actors (
  name,
  gender,
  birthday
) VALUES 
  ($1, $2, $3) RETURNING id, name, gender, birthday
`

type CreateActorParams struct {
	Name     string    `json:"name"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

func (q *Queries) CreateActor(ctx context.Context, arg CreateActorParams) (Actor, error) {
	row := q.db.QueryRowContext(ctx, createActor, arg.Name, arg.Gender, arg.Birthday)
	var i Actor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Gender,
		&i.Birthday,
	)
	return i, err
}

const deleteActor = `-- name: DeleteActor :exec
DELETE FROM actors
WHERE id = $1
`

func (q *Queries) DeleteActor(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteActor, id)
	return err
}

const getActorMoviesList = `-- name: GetActorMoviesList :many
SELECT
    a.id AS actor_id,
    a.name AS actor_name,
    m.id AS movie_id,
    m.name AS movie_name
FROM
    actors a
LEFT JOIN
    movie_actors ma ON a.id = ma.actor_id
LEFT JOIN
    movies m ON ma.movie_id = m.id
ORDER BY
    a.id, m.id
`

type GetActorMoviesListRow struct {
	ActorID   int32          `json:"actor_id"`
	ActorName string         `json:"actor_name"`
	MovieID   sql.NullInt32  `json:"movie_id"`
	MovieName sql.NullString `json:"movie_name"`
}

func (q *Queries) GetActorMoviesList(ctx context.Context) ([]GetActorMoviesListRow, error) {
	rows, err := q.db.QueryContext(ctx, getActorMoviesList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetActorMoviesListRow{}
	for rows.Next() {
		var i GetActorMoviesListRow
		if err := rows.Scan(
			&i.ActorID,
			&i.ActorName,
			&i.MovieID,
			&i.MovieName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateActor = `-- name: UpdateActor :one
UPDATE actors
SET name = $2,
  gender = $3,
  birthday = $4
WHERE id = $1
RETURNING id, name, gender, birthday
`

type UpdateActorParams struct {
	ID       int32     `json:"id"`
	Name     string    `json:"name"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

func (q *Queries) UpdateActor(ctx context.Context, arg UpdateActorParams) (Actor, error) {
	row := q.db.QueryRowContext(ctx, updateActor,
		arg.ID,
		arg.Name,
		arg.Gender,
		arg.Birthday,
	)
	var i Actor
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Gender,
		&i.Birthday,
	)
	return i, err
}

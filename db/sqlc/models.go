// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"
)

type Actor struct {
	ID       int32     `json:"id"`
	Name     string    `json:"name"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

type Movie struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      string    `json:"rating"`
}

type MovieActor struct {
	MovieID int32 `json:"movie_id"`
	ActorID int32 `json:"actor_id"`
}

type User struct {
	ID                int32     `json:"id"`
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	Role              string    `json:"role"`
	CreatedAt         time.Time `json:"created_at"`
}

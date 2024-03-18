package api

import (
	"database/sql"
	"net/http"
	"time"
	db "vk-film/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// createMovieRequest represents a request to create a movie, ONLY FOR ADMINS.
// swagger:parameters createMovieRequest
type createMovieRequest struct {
	// The name of the movie.
	// Required: true
	// example: Inception
	Name string `json:"name" binding:"required"`

	// The description of the movie.
	// Required: true
	// example: A mind-bending action thriller
	Description string `json:"description" binding:"required"`

	// The release date of the movie.
	// Required: true
	// format: date
	// example: "2010-07-16"
	ReleaseDate time.Time `json:"release_date" binding:"required"`

	// The rating of the movie.
	// Required: true
	// example: 8.8
	Rating string `json:"rating" binding:"required"`
}

// movieResponse represents the response for a movie.
// swagger:response movieResponse
type movieResponse struct {
	// The ID of the movie.
	// Example: 123
	// required: true
	ID int32 `json:"id"`

	// The name of the movie.
	// Example: The Matrix
	// required: true
	Name string `json:"name"`

	// The description of the movie.
	// Example: A computer hacker learns from mysterious rebels about the true nature of his reality and his role in the war against its controllers.
	// required: true
	Description string `json:"description"`

	// The release date of the movie.
	// Example: 1999-03-31
	// required: true
	ReleaseDate time.Time `json:"release_date"`

	// The rating of the movie.
	// Example: 8.7
	// required: true
	Rating string `json:"rating"`
}

// newMovieResponse creates a new Movie Response from a db.Movie.
func newMovieResponse(movie db.Movie) movieResponse {
	return movieResponse{
		ID:          movie.ID,
		Name:        movie.Name,
		Description: movie.Description,
		ReleaseDate: movie.ReleaseDate,
		Rating:      movie.Rating,
	}
}

// createMovie creates a new movie.
// swagger:operation POST /movie/create movies createMovie
//
// Creates a new movie with the provided details.
//
// ---
// produces:
// - application/json
// parameters:
//   - name: body
//     in: body
//     description: Request body containing the movie information.
//     required: true
//     schema:
//     "$ref": "#/definitions/createMovieRequest"
//
// responses:
//
//	'200':
//	  description: Successfully created the movie.
//	  schema:
//	    "$ref": "#/definitions/movieResponse"
//	'400':
//	  description: Bad request. The request body is missing or invalid.
//	'403':
//	  description: Forbidden. Only admins have permission to update movies.
//	'404':
//	  description: Not found. The movie with the provided ID does not exist.
//	'500':
//	  description: Internal server error. Something went wrong while processing the request.
func (server *Server) createMovie(ctx *gin.Context) {
	var req createMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err := checkAdminPermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	arg := db.CreateMovieParams{
		Name:        req.Name,
		Description: req.Description,
		ReleaseDate: req.ReleaseDate,
		Rating:      req.Rating,
	}
	movie, err := server.store.CreateMovie(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newMovieResponse(movie)
	ctx.JSON(http.StatusOK, rsp)
}

// updateMovieRequest represents the request body for updating a movie.
// swagger:parameters updateMovie
type updateMovieRequest struct {
	// ID of the movie to be updated.
	// required: true
	// in: body
	ID int32 `json:"id" binding:"required"`

	// New name of the movie.
	// in: body
	Name string `json:"name"`

	// New description of the movie.
	// in: body
	Description string `json:"description"`

	// New rating of the movie.
	// in: body
	Rating string `json:"rating"`

	// New release date of the movie.
	// in: body
	ReleaseDate time.Time `json:"release_date"`
}

// updateMovie updates a movie based on the provided request body.
// swagger:route PATCH /movie/update movies updateMovie
// responses:
//
//	'200': movieResponse
//	'400':
//	  description: Bad request. The request body is missing or invalid.
//	'403':
//	  description: Forbidden. Only admins have permission to update movies.
//	'404':
//	  description: Not found. The movie with the provided ID does not exist.
//	'500':
//	  description: Internal server error. Something went wrong while processing the request.
func (server *Server) updateMovie(ctx *gin.Context) {
	var req updateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := checkAdminPermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	arg := db.UpdateMovieParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Rating:      req.Rating,
		ReleaseDate: req.ReleaseDate,
	}
	movie, err := server.store.UpdateMovie(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "no_data_found":
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movie)

}

// deleteMovieRequest represents the request payload for deleting a movie.
// swagger:model
type deleteMovieRequest struct {
	// The ID of the movie to delete.
	//
	// required: true
	// example: 123
	ID int32 `form:"id" binding:"required"`
}

// deleteMovie deletes a movie based on the provided ID.
// swagger:operation DELETE /movie/delete movies deleteMovie
//
// Delete a movie with the provided ID.
//
// ---
// produces:
// - application/json
// parameters:
//   - name: body
//     in: body
//     description: Request body containing the ID of the movie to delete.
//     required: true
//     schema:
//     "$ref": "#/definitions/deleteMovieRequest"
//
// responses:
//
//	'200':
//	  description: Successfully deleted the movie.
//	  schema:
//	    type: integer
//	    format: int32
//	'400':
//	  description: Bad request. The request body is missing or invalid.
//	'403':
//	  description: Forbidden. The user does not have permission to delete movies.
//	'404':
//	  description: Not found. The movie with the provided ID does not exist.
func (server *Server) deleteMovie(ctx *gin.Context) {
	var req deleteMovieRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := checkAdminPermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	err = server.store.DeleteMovie(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, req.ID)
}

// moviesSortedByRating retrieves movies sorted by rating.
// swagger:operation GET /movies movies moviesSortedByRating
//
// Retrieves movies sorted by rating.
// ---
// responses:
//
//	200:
//	  description: Successfully retrieved movies sorted by rating.
//	  schema:
//	    type: array
//	    items:
//	      $ref: "#/definitions/MovieResponse"
//	500:
//	  description: Internal server error.
func (server *Server) moviesSortedByRating(ctx *gin.Context) {
	movies, err := server.store.GetMoviesSortedByRating(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

// moviesSortedByName retrieves movies sorted by name.
// swagger:operation GET /movies/by-name movies moviesSortedByName
//
// Retrieves movies sorted by name.
// ---
// responses:
//
//	200:
//	  description: Successfully retrieved movies sorted by name.
//	  schema:
//	    type: array
//	    items:
//	      $ref: "#/definitions/MovieResponse"
//	500:
//	  description: Internal server error.
func (server *Server) moviesSortedByName(ctx *gin.Context) {
	movies, err := server.store.GetMoviesSortedByName(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

// moviesSortedByReleaseDate retrieves movies sorted by release date.
// swagger:operation GET /movies/by-date movies moviesSortedByReleaseDate
//
// Retrieves movies sorted by release date.
// ---
// responses:
//
//	200:
//	  description: Successfully retrieved movies sorted by release date.
//	  schema:
//	    type: array
//	    items:
//	      $ref: "#/definitions/MovieResponse"
//	500:
//	  description: Internal server error.
func (server *Server) moviesSortedByReleaseDate(ctx *gin.Context) {
	movies, err := server.store.GetMoviesByReleaseDate(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

type movieFragmentRequest struct {
	Fragment string `form:"name" binding:"required"`
}

// moviesByNameFragment retrieves movies by name fragment.
// swagger:operation GET /movies/by-name-fragment movies moviesByNameFragment
//
// Retrieves movies by name fragment.
//
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
//   - name: body
//     in: body
//     description: The name fragment to search for.
//     required: true
//     schema:
//     $ref: "#/definitions/MovieFragmentRequest"
//
// responses:
//
//	200:
//	  description: Successfully retrieved movies by name fragment.
//	  schema:
//	    type: array
//	    items:
//	      $ref: "#/definitions/MovieResponse"
//	400:
//	  description: Bad request. Invalid JSON payload.
func (server *Server) moviesByNameFragment(ctx *gin.Context) {
	var req movieFragmentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fragment := req.Fragment
	sqlStr := sql.NullString{
		String: fragment,
		Valid:  true,
	}
	movies, err := server.store.GetMoviesByNameFragment(ctx, sqlStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

type actorFragmentRequest struct {
	Fragment string `form:"actor" binding:"required"`
}

// moviesByActorFragment retrieves movies by actor fragment.
// swagger:operation GET /movies/by-actor-fragment movies moviesByActorFragment
//
// Retrieves movies by actor fragment.
//
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
//   - name: body
//     in: body
//     description: The name fragment to search for.
//     required: true
//     schema:
//     $ref: "#/definitions/ActorFragmentRequest"
//
// responses:
//
//	200:
//	  description: Successful response with an array of movie objects
//	  schema:
//	    type: array
//	    items:
//	      $ref: "#/definitions/MovieResponse"
//	400:
//	  description: Bad request response
//	  schema:
//	    $ref: "#/definitions/ErrorResponse"
//	500:
//	  description: Internal server error response
//	  schema:
//	    $ref: "#/definitions/ErrorResponse"
func (server *Server) moviesByActorFragment(ctx *gin.Context) {
	var req actorFragmentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fragment := req.Fragment
	sqlStr := sql.NullString{
		String: fragment,
		Valid:  true,
	}
	movies, err := server.store.GetMoviesByActorFragment(ctx, sqlStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

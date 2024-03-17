package api

import (
	"database/sql"
	"net/http"
	"time"
	db "vk-film/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createMovieRequest struct {
	Name        string    `json:"name"         binding:"required"`
	Description string    `json:"description"  binding:"required"`
	ReleaseDate time.Time `json:"release_date" binding:"required"`
	Rating      string    `json:"rating"       binding:"required"`
}

type movieResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      string    `json:"rating"`
}

func newMovieResponse(movie db.Movie) movieResponse {
	return movieResponse{
		ID:          movie.ID,
		Name:        movie.Name,
		Description: movie.Description,
		ReleaseDate: movie.ReleaseDate,
		Rating:      movie.Rating,
	}
}

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

type updateMovieRequest struct {
	ID          int32     `json:"id"           binding:"required"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rating      string    `json:"rating"`
	ReleaseDate time.Time `json:"release_date"`
}

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

type deleteMovieRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func (server *Server) deleteMovie(ctx *gin.Context) {
	var req deleteMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

func (server *Server) moviesSortedByRating(ctx *gin.Context) {
	movies, err := server.store.GetMoviesSortedByRating(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

func (server *Server) moviesSortedByName(ctx *gin.Context) {
	movies, err := server.store.GetMoviesSortedByName(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

func (server *Server) moviesSortedByReleaseDate(ctx *gin.Context) {
	movies, err := server.store.GetMoviesByReleaseDate(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

type movieFragmentRequest struct {
	Fragment string `json:"name" binding:"required"`
}

func (server *Server) moviesByNameFragment(ctx *gin.Context) {
	var req movieFragmentRequest
	if err := ctx.BindJSON(&req); err != nil {
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
	Fragment string `json:"actor" binding:"required"`
}

func (server *Server) moviesByActorFragment(ctx *gin.Context) {
	var req actorFragmentRequest
	if err := ctx.BindJSON(&req); err != nil {
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

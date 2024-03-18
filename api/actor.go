package api

import (
	"net/http"
	"time"
	db "vk-film/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// createActorRequest represents the request body for creating an actor.
// swagger:parameters createActor
type createActorRequest struct {
	// The name of the actor.
	// Required: true
	Name string `json:"name" binding:"required"`

	// The gender of the actor.
	// Required: true
	Gender string `json:"gender" binding:"required"`

	// The birthday of the actor.
	// Required: true
	Birthday time.Time `json:"birthday" binding:"required"`
}

// actorResponse represents the response body for an actor.
// swagger:response actorResponse
type actorResponse struct {
	// The ID of the actor.
	// Example: 1
	ID int32 `json:"id"`

	// The name of the actor.
	// Example: Leonardo Di Caprio
	Name string `json:"name"`

	// The gender of the actor, can be: ["male", "female", "other"].
	// Example: male
	Gender string `json:"gender"`

	// The birthday of the actor.
	// Example: 2000-01-01T00:00:00Z
	Birthday time.Time `json:"birthday"`
}

// newActorResponse creates a new actorResponse from a db.Actor.
func newActorResponse(actor db.Actor) actorResponse {
	return actorResponse{
		ID:       actor.ID,
		Name:     actor.Name,
		Gender:   actor.Gender,
		Birthday: actor.Birthday,
	}
}

// createActor creates a new actor.
// swagger:route POST /actors actors createActor
// Creates a new actor.
// responses:
//
//	200: actorResponse
//	400: errorResponse
//	401: errorResponse
//	500: errorResponse
func (server *Server) createActor(ctx *gin.Context) {
	var req createActorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err := checkAdminPermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	arg := db.CreateActorParams{
		Name:     req.Name,
		Gender:   req.Gender,
		Birthday: req.Birthday,
	}
	actor, err := server.store.CreateActor(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newActorResponse(actor)
	ctx.JSON(http.StatusOK, rsp)
}

// updateActorRequest represents the request body for updating an actor.
// swagger:parameters updateActor
type updateActorRequest struct {
	// The ID of the actor to update.
	// Required: true
	ID int32 `json:"id" binding:"required"`

	// The name of the actor.
	// Example: John Doe
	Name string `json:"name"`

	// The gender of the actor.
	// Example: male
	Gender string `json:"gender"`

	// The birthday of the actor.
	// Example: 2000-01-01T00:00:00Z
	Birthday time.Time `json:"birthday"`
}

// updateActor updates an existing actor.
// swagger:route PUT /actors actors updateActor
// Updates an existing actor.
// responses:
//
//		'200': actorResponse
//		'400':
//		  description: Bad request. The request body is missing or invalid.
//		'403':
//		  description: Forbidden. Only admins have permission to update actors.
//		'404':
//		  description: Not found. The actor with the provided ID does not exist.
//	 '500':
//		  description: Internal server error. Something went wrong while processing the request.
func (server *Server) updateActor(ctx *gin.Context) {
	var req updateActorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := checkAdminPermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	arg := db.UpdateActorParams{
		ID:       req.ID,
		Name:     req.Name,
		Gender:   req.Gender,
		Birthday: req.Birthday,
	}
	actor, err := server.store.UpdateActor(ctx, arg)
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
	ctx.JSON(http.StatusOK, actor)
}

// deleteActorRequest represents the request body for deleting an actor.
// swagger:parameters deleteActor
type deleteActorRequest struct {
	// The ID of the actor to delete.
	// in: body
	// required: true
	ID int32 `form:"id" binding:"required"`
}

// deleteActor deletes an actor by ID.
// swagger:route DELETE /actors actors deleteActor
// Deletes an actor by ID.
// responses:
//
//	204: noContentResponse
//	400: errorResponse
//	403: errorResponse
//	404: errorResponse
func (server *Server) deleteActor(ctx *gin.Context) {
	var req deleteActorRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := checkAdminPermissions(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	err = server.store.DeleteActor(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusNoContent, req.ID)
}

func (server *Server) actorsWithMovies(ctx *gin.Context) {
	actorsWithMovies, err := server.store.GetActorMoviesList(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, actorsWithMovies)
}

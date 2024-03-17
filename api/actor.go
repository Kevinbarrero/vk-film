package api

import (
	"net/http"
	"time"
	db "vk-film/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createActorRequest struct {
	Name     string    `json:"name"     binding:"required"`
	Gender   string    `json:"gender"   binding:"required"`
	Birthday time.Time `json:"birthday" binding:"required"`
}

type actorResponse struct {
	ID       int32     `json:"id"`
	Name     string    `json:"name"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

func newActorResponse(actor db.Actor) actorResponse {
	return actorResponse{
		ID:       actor.ID,
		Name:     actor.Name,
		Gender:   actor.Gender,
		Birthday: actor.Birthday,
	}
}

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

type updateActorRequest struct {
	ID       int32     `json:"id"       binding:"required"`
	Name     string    `json:"name"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

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

type deleteActorRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func (server *Server) deleteActor(ctx *gin.Context) {
	var req deleteActorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
	ctx.JSON(http.StatusOK, req.ID)
}

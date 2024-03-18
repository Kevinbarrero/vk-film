package api

import (
	"errors"
	"fmt"
	db "vk-film/db/sqlc"
	"vk-film/token"
	"vk-film/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// movie routes
	authRoutes.POST("/movie/create", server.createMovie)
	authRoutes.PATCH("/movie/update", server.updateMovie)
	authRoutes.DELETE("/movie/delete", server.deleteMovie)
	authRoutes.GET("/movies", server.moviesSortedByRating)
	authRoutes.GET("/movies/by-name", server.moviesSortedByName)
	authRoutes.GET("/movies/by-date", server.moviesSortedByReleaseDate)
	authRoutes.GET("/movies/by-name-fragment", server.moviesByNameFragment)
	authRoutes.GET("/movies/by-actor-fragment", server.moviesByActorFragment)

	// actor routes

	authRoutes.POST("/actor/create", server.createActor)
	authRoutes.PATCH("/actor/update", server.updateActor)
	authRoutes.DELETE("/actor/delete", server.deleteActor)
	authRoutes.GET("/actors-movies", server.actorsWithMovies)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func checkAdminPermissions(ctx *gin.Context) error {
	authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
	if authPayload.Role != "administrator" {
		err := errors.New("Insufficient permissions, only admin can modify data")
		return err
	}
	return nil
}

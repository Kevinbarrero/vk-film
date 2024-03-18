package api

import (
	"database/sql"
	"net/http"
	"time"
	db "vk-film/db/sqlc"
	"vk-film/util"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// createUserRequest represents the request body for creating a new user, the user's role cannot be setted through api, by default is 'client', if you want to check administrator endpoints you can create an user directly in the database, otherwise you can use the default administrator [username: 'admin', password: 'qwerty'].
// swagger:parameters createUser
type createUserRequest struct {
	// The username of the user.
	// Required: true
	// Example: vk-qwerty
	Username string `json:"username" binding:"required,alphanum"`

	// The password of the user.
	// Required: true
	// Example: password123
	Password string `json:"password" binding:"required,alphanum"`
}

// userResponse represents the response body for a user.
// swagger:response userResponse
type userResponse struct {
	// The username of the user.
	// Example: vk-user
	Username string `json:"username"`

	// The timestamp when the password was last changed.
	// Example: "2022-03-17T10:00:00Z"
	PasswordChangedAt time.Time `json:"password_changed_at"`

	// The timestamp when the user was created.
	// Example: "2022-03-17T09:00:00Z"
	CreatedAt time.Time `json:"created_at"`

	// The role of the user.
	// Example: administrator
	Role string `json:"role"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		Role:              user.Role,
	}
}

// createUser creates a new user.
// swagger:route POST /users users createUser
// Creates a new user.
//
// Consumes:
// - application/json
//
// Parameters:
//   - In: body
//     name: user
//     description: User to create.
//     schema:
//     type: object
//     required:
//   - username
//   - password
//     properties:
//     username:
//     type: string
//     description: The username of the user.
//     password:
//     type: string
//
// Responses:
//
//	200: userResponse
//	400: errorResponse
//	403: errorResponse
//	500: errorResponse
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

// loginUserRequest represents the request body for logging in a user.
// swagger:parameters loginUser
// in: body
type loginUserRequest struct {
	// The username of the user.
	// Required: true
	// Example: john_doe
	Username string `json:"username" binding:"required,alphanum"`

	// The password of the user.
	// Required: true
	// Minimum: 6
	// Example: password123
	Password string `json:"password" binding:"required,min=6"`
}

// loginUserResponse represents the response body for logging in a user.
// swagger:response loginUser
type loginUserResponse struct {
	// The access token for the user.
	// Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im...
	AccessToken string `json:"access_token"`

	// The user information.
	// Example: {"username":"admin","password_changed_at":"2022-03-17T10:00:00Z","created_at":"2022-03-17T09:00:00Z","role":"admin"}
	User userResponse `json:"user"`
}

// loginUser logs in a user based on the provided request body.
// swagger:route POST /users/login users loginUser
//
// Logs in a user.
//
// responses:
//
//	'200': loginUser
//	'400':
//	  description: Bad request. The request body is missing or invalid.
//	'401':
//	  description: Unauthorized. Username or password is incorrect.
//	'500':
//	  description: Internal server error. Something went wrong while processing the request.
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

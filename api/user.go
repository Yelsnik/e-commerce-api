package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type signUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" `
	Email     string    `json:"email" `
	Role      string    `json:"role" `
	CreatedAt time.Time `json:"created_at"`
}

func newSignUpResponse(user db.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) signUp(ctx *gin.Context) {
	// unmarshal the request
	var req signUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash the password
	hashedPassword, err := util.HashPassword(req.Password)

	arg := db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: hashedPassword,
	}

	// create the user in the db
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response := newSignUpResponse(user)

	success(ctx, response)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func newLoginResponse(user db.User, accessToken string) loginUserResponse {
	return loginUserResponse{
		AccessToken: accessToken,
		User: userResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.ComparePassword(req.Password, user.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.ID,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newLoginResponse(user, accessToken)

	success(ctx, response)

}

package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/mail"
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

	accessToken, _, err := server.tokenMaker.CreateToken(
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

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (server *Server) forgotPassword(ctx *gin.Context) {
	var req forgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	resetToken, payload, err := server.tokenMaker.CreateToken(user.ID, user.Role, server.config.PasswordResetTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: payload.ExpiredAt,
	}

	_, err = server.store.CreatePasswordResetToken(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resetLink := fmt.Sprintf("localhost/v1/forgot-password?token=%s", resetToken)

	sender := mail.NewGmailSender(server.config.EmailSenderName, server.config.EmailSenderAddress, server.config.EmailSenderPassword)

	subject := "Reset your password"

	content := fmt.Sprintf(
		`
	<h1> Reset password link </h1>
	<p> click this link to reset password <a href=%s>link<a/> </p>
	`, resetLink)

	to := []string{user.Email}

	attachFiles := []string{}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "password reset link sent succesfully"})
}

type resetPasswordRequest struct {
	ResetToken string `json:"reset_token"`
	Password   string `json:"password"`
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.tokenMaker.VerifyToken(req.ResetToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	passwordResetToken, err := server.store.GetPasswordResetTokenByToken(ctx, req.ResetToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:       passwordResetToken.UserID,
		Password: hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully updated password", "data": user})
}

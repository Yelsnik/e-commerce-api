package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Yelsnik/e-commerce-api/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func roleBasedMiddleware(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get the payload
		payload, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			err := errors.New("you are not logged in")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Assert payload type
		authPayload, ok := payload.(*token.Payload)
		if !ok {
			err := errors.New("invalid authorization payload")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Check if the user's role matches the required role
		var merchant string
		var admin string
		for i := 0; i < len(roles); i++ {
			if roles[i] == "merchant" {
				merchant = roles[i]
			}
			if roles[i] == "admin" {
				admin = roles[i]
			}
		}

		if authPayload.Role == merchant {
			ctx.Next()
			return
		}

		if authPayload.Role == admin {
			ctx.Next()
			return
		}

		err := errors.New("you are not authorised to perform this action")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return

	}
}

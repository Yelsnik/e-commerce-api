package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Yelsnik/e-commerce-api/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func addAuthentication(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	user_id uuid.UUID, role string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(user_id, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

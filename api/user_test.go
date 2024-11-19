package api

import (
	"testing"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:       util.Test(),
		Name:     util.RandomString(6),
		Email:    util.RandomEmail(),
		Role:     "admin",
		Password: hashedPassword,
	}
	return
}

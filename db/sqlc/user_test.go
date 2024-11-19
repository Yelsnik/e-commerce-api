package db

import (
	//"fmt"
	"context"
	"testing"
	"time"

	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func createNewUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Name:     util.RandomString(6),
		Email:    util.RandomEmail(),
		Role:     util.RandomRole(),
		Password: hashedPassword,
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Role, user.Role)
	require.NotZero(t, user.CreatedAt)
	//require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createNewUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createNewUser(t)

	user, err := testStore.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user1.ID, user.ID)
	require.Equal(t, user1.Password, user.Password)
	require.Equal(t, user1.Name, user.Name)
	require.Equal(t, user1.Email, user.Email)
	require.Equal(t, user1.Role, user.Role)
	// require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user.CreatedAt, time.Second)
}

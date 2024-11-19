package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createNewCart(t *testing.T) Cart {
	user := createNewUser(t)

	args := CreateCartsParams{
		UserID:     user.ID,
		TotalPrice: float64(0),
	}

	cart, err := testStore.CreateCarts(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, cart)
	require.Equal(t, args.UserID, cart.UserID)
	require.Equal(t, args.TotalPrice, cart.TotalPrice)

	return cart
}

func TestCreateCart(t *testing.T) {
	createNewCart(t)
}

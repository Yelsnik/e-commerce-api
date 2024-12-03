package db

import (
	"context"
	"testing"

	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func createNewOrder(t *testing.T) Order {

	user := createNewUser(t)

	arg := CreateOrdersParams{
		UserName:        user.Name,
		UserID:          user.ID,
		TotalPrice:      float64(util.RandomMoney()),
		DeliveryAddress: util.RandomString(8),
		Country:         util.RandomCountry(),
		Status:          "processing",
	}

	order, err := testStore.CreateOrders(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, order)
	require.Equal(t, arg.UserName, order.UserName)
	require.Equal(t, arg.UserID, order.UserID)
	require.Equal(t, arg.TotalPrice, order.TotalPrice)
	require.Equal(t, arg.DeliveryAddress, order.DeliveryAddress)
	require.Equal(t, arg.Country, order.Country)
	require.Equal(t, arg.Status, order.Status)
	require.NotZero(t, order.CreatedAt)

	return order
}

func TestCreateNewOrder(t *testing.T) {
	createNewOrder(t)
}

package db

import (
	"context"
	"testing"

	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func TestCreateOrderTx(t *testing.T) {
	var cartItems []Cartitem

	n := 5

	for i := 0; i < n; i++ {
		cartItems = append(cartItems, createNewCartItem(t))

	}

	for _, cartItem := range cartItems {

		cart, errC := testStore.GetCarts(context.Background(), cartItem.Cart)
		require.NoError(t, errC)

		product, errP := testStore.GetProducts(context.Background(), cartItem.Product)
		require.NoError(t, errP)

		user, errU := testStore.GetUser(context.Background(), cart.UserID)
		require.NoError(t, errU)

		arg := OrderTxParams{
			UserName:        user.Name,
			UserID:          user.ID,
			TotalPrice:      cartItem.SubTotal,
			DeliveryAddress: util.RandomString(8),
			Country:         util.RandomCountry(),
		}

		errs := make(chan error, 1)
		results := make(chan OrderTxResult, 1)

		go func() {
			result, err := testStore.CreateOrderTx(context.Background(), cart.ID, cartItem.ID, arg)

			errs <- err
			results <- result
		}()

		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		order := result.Order
		require.NotEmpty(t, order)
		require.Equal(t, arg.UserName, order.UserName)
		require.Equal(t, arg.UserID, order.UserID)
		require.Equal(t, arg.TotalPrice, order.TotalPrice)
		require.Equal(t, arg.DeliveryAddress, order.DeliveryAddress)
		require.Equal(t, arg.Country, order.Country)
		require.NotZero(t, order.CreatedAt)

		orderItem := result.OrderItem
		require.NotEmpty(t, orderItem)
		require.Equal(t, product.ProductName, orderItem.ItemName)
		require.Equal(t, cartItem.SubTotal, orderItem.ItemSubTotal)
		require.Equal(t, cartItem.Quantity, orderItem.Quantity)
		require.Equal(t, product.ID, orderItem.ItemID)
		require.Equal(t, order.ID, orderItem.OrderID)
		require.NotZero(t, orderItem.CreatedAt)

		//	require.Empty(t, result.CartItem)
	}

}

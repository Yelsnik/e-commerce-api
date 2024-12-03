package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createNewOrderItem(t *testing.T) Orderitem {
	product := createNewProduct(t)
	order := createNewOrder(t)

	quantity := 4

	subTotal := float64(quantity) * product.Price

	arg := CreateOrderitemsParams{
		ItemName:     product.ProductName,
		ItemSubTotal: subTotal,
		Quantity:     int64(quantity),
		ItemID:       product.ID,
		OrderID:      order.ID,
	}

	orderItem, err := testStore.CreateOrderitems(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, orderItem)
	require.Equal(t, arg.ItemName, orderItem.ItemName)
	require.Equal(t, arg.ItemSubTotal, orderItem.ItemSubTotal)
	require.Equal(t, arg.Quantity, orderItem.Quantity)
	require.Equal(t, arg.ItemID, orderItem.ItemID)
	require.Equal(t, arg.OrderID, orderItem.OrderID)
	require.NotZero(t, order.CreatedAt)

	return orderItem

}

func TestCreateNewOrderItem(t *testing.T) {
	createNewOrderItem(t)
}

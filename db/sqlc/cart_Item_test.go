package db

import (
	"context"
	"testing"

	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func createNewCartItem(t *testing.T) Cartitem {
	cart := createNewCart(t)
	product := createNewProduct(t)

	q := util.RandomInt(1, 10)
	p := float64(util.RandomMoney())
	subTotal := float64(q) * p

	arg := CreateCartitemsParams{
		Cart:     cart.ID,
		Product:  product.ID,
		Quantity: q,
		Price:    p,
		Currency: util.RandomCurrency(),
		SubTotal: subTotal,
	}

	cartItem, err := testStore.CreateCartitems(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, cartItem)

	require.Equal(t, arg.Cart, cartItem.Cart)
	require.Equal(t, arg.Product, cartItem.Product)
	require.Equal(t, arg.Quantity, cartItem.Quantity)
	require.Equal(t, arg.Price, cartItem.Price)
	require.Equal(t, arg.SubTotal, cartItem.SubTotal)
	require.NotEmpty(t, cartItem.CreatedAt)
	require.NotZero(t, cartItem.CreatedAt)

	return cartItem
}

func TestCreateCartItem(t *testing.T) {
	createNewCartItem(t)
}

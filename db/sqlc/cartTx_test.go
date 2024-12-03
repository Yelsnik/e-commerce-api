package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddToCartTx(t *testing.T) {
	var product []Product
	for i := 0; i < 3; i++ {
		var arr []Product
		product = append(arr, createNewProduct(t))
	}

	cart := createNewCart(t)

	n := 10

	for i := 0; i < len(product); i++ {

		q := int64(3)

		subTotal := float64(q) * product[i].Price
		arg := CreateCartitemsParams{
			Cart:     cart.ID,
			Product:  product[i].ID,
			Quantity: q,
			Price:    product[i].Price,
			Currency: product[i].Currency,
			SubTotal: subTotal,
		}

		errs := make(chan error, 5)
		results := make(chan CartTxResult, 5)

		for i := 0; i < n; i++ {
			txName := fmt.Sprintf("tx %d", i+1)

			go func() {
				ctx := context.WithValue(context.Background(), txKey, txName)
				result, err := testStore.AddToCartTx(ctx, arg)

				results <- result
				errs <- err
			}()
		}

		// check result
		for i := 0; i < n; i++ {
			err := <-errs
			require.NoError(t, err)

			result := <-results
			require.NotEmpty(t, result)

			cartResult := result.Cart
			require.NotEmpty(t, cartResult)
			require.Equal(t, cart.ID, cartResult.ID)

			cartItemResult := result.CartItem
			require.NotEmpty(t, cartItemResult)

		}

	}

}

func TestUpdateCartTx(t *testing.T) {
	cartItem := createNewCartItem(t)

	n := 3
	updatedQuantity := int64(10)
	updatedSubtotal := float64(updatedQuantity) * cartItem.Price

	arg := UpdateCartitemsParams{
		ID:       cartItem.ID,
		Quantity: int64(updatedQuantity),
		SubTotal: updatedSubtotal,
	}

	errs := make(chan error, 5)
	results := make(chan CartTxResult, 5)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := testStore.UpdateCartTx(ctx, cartItem.ID, arg)

			errs <- err
			results <- result
		}()
	}

	// check result and err
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		cartResult := result.Cart
		require.NotEmpty(t, cartResult)

		cartItemResult := result.CartItem
		require.NotEmpty(t, cartItemResult)
		require.Equal(t, updatedQuantity, cartItemResult.Quantity)
		require.Equal(t, updatedSubtotal, cartItemResult.SubTotal)
	}

}

func TestRemoveCartTx(t *testing.T) {
	cartItem := createNewCartItem(t)

	n := 1

	errs := make(chan error, 5)
	results := make(chan RemoveCartTxResult, 5)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := testStore.RemoveCartTx(ctx, cartItem.ID, cartItem.Cart)

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results

		cart := result.Cart
		require.NotEmpty(t, cart)

	}

}

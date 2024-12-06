package db

import (
	"context"
	"testing"

	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func createNewPayment(t *testing.T, user User) Payment {

	arg := CreatePaymentsParams{
		ID:       util.Test().String(),
		Amount:   float64(util.RandomMoney()),
		Currency: util.RandomCurrency(),
		Status:   "processing",
		UserID:   user.ID,
	}

	payment, err := testStore.CreatePayments(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, arg.ID, payment.ID)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.Currency, payment.Currency)
	require.Equal(t, arg.Status, payment.Status)
	require.Equal(t, arg.UserID, payment.UserID)

	return payment
}

func TestCreatePayment(t *testing.T) {
	user := createNewUser(t)

	createNewPayment(t, user)
}

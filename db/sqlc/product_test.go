package db

import (
	"context"
	"testing"

	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func createNewProduct(t *testing.T) Product {
	user := createNewUser(t)
	arg := CreateProductsParams{
		Category:     util.RandomCategory(),
		ProductName:  util.RandomString(6),
		Description:  util.RandomString(7),
		Brand:        util.NewNullString(util.RandomString(6)),
		CountInStock: util.RandomInt(0, 9),
		Price:        float64(util.RandomMoney()),
		Currency:     util.RandomCurrency(),
		Rating:       util.NewNullInt(util.RandomInt(0, 5)),
		UserID:       user.ID,
	}

	product, err := testStore.CreateProducts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)

	require.Equal(t, arg.Category, product.Category)
	require.Equal(t, arg.ProductName, product.ProductName)
	require.Equal(t, arg.Description, product.Description)
	require.Equal(t, arg.Brand, product.Brand)
	require.Equal(t, arg.CountInStock, product.CountInStock)
	require.Equal(t, arg.Price, product.Price)
	require.Equal(t, arg.Currency, product.Currency)
	require.Equal(t, arg.Rating, product.Rating)
	require.Equal(t, arg.UserID, product.UserID)

	return product
}

func TestCreateProduct(t *testing.T) {
	createNewProduct(t)
}

func TestGetProduct(t *testing.T) {
	product := createNewProduct(t)

	p, err := testStore.GetProducts(context.Background(), product.ID)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	require.Equal(t, product.Category, p.Category)
	require.Equal(t, product.ProductName, p.ProductName)
	require.Equal(t, product.Description, p.Description)
	require.Equal(t, product.Brand, p.Brand)
	require.Equal(t, product.CountInStock, p.CountInStock)
	require.Equal(t, product.Price, p.Price)
	require.Equal(t, product.Currency, p.Currency)
	require.Equal(t, product.Rating, p.Rating)
}

func TestListProduct(t *testing.T) {
	arg := ListProductsParams{
		Limit:  5,
		Offset: (1 - 1) * 5,
	}

	products, err := testStore.ListProducts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, products)

	for _, product := range products {
		require.NotEmpty(t, product)
	}
}

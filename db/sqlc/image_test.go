package db

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func createNewImage(t *testing.T, product Product) Image {

	var image Image
	file := "image_test.jpg"

	contents, err := os.ReadFile(file)
	require.NoError(t, err)

	arg := CreateImagesParams{
		ImageName: "image_test.jpg",
		Data:      contents,
		Product:   product.ID,
	}

	image, err = testStore.CreateImages(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, arg.ImageName, image.ImageName)
	require.Equal(t, arg.Data, image.Data)
	require.Equal(t, arg.Product, image.Product)
	require.NotZero(t, image.CreatedAt)

	return image
}

func TestCreateNewImage(t *testing.T) {
	product := createNewProduct(t)
	createNewImage(t, product)
}

func TestGetImage(t *testing.T) {

	product := createNewProduct(t)
	createNewImage(t, product)

	image, err := testStore.GetImages(context.Background(), product.ID)
	require.NoError(t, err)

	require.Equal(t, product.ID, image.Product)
	require.NotEmpty(t, image.Data)

}

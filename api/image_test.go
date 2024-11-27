package api

import (
	"os"
	"testing"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func randomImage(t *testing.T, product db.Product, image string) db.Image {
	file := "image_test.jpg"

	contents, err := os.ReadFile(file)
	require.NoError(t, err)

	return db.Image{
		ID:        util.RandomInt(1, 1000),
		ImageName: image,
		Data:      contents,
		Product:   product.ID,
	}
}

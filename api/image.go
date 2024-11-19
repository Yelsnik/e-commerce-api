package api

import (
	"bytes"

	"io"
	"mime/multipart"
	"net/http"

	//	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const image = "image"

func fileToBytes(file multipart.File) ([]byte, error) {
	// Use a bytes.Buffer to hold the file's data
	var buf bytes.Buffer
	// Copy the file's data into the buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}
	// Return the buffer's contents as a byte slice
	return buf.Bytes(), nil
}

type uploadImageRequest struct {
	ProductID string `uri:"pid" binding:"required"`
}

type uploadImageResponse struct {
	ImageName string    `json:"image_name"`
	Data      []byte    `json:"data"`
	Product   uuid.UUID `json:"product"`
}

func newUploadImageResponse(image db.Image) uploadImageResponse {
	return uploadImageResponse{
		ImageName: image.ImageName,
		Data:      image.Data,
		Product:   image.Product,
	}
}

func (server *Server) uploadImage(ctx *gin.Context) {
	var req uploadImageRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// parse form
	fileHeader, err := ctx.FormFile(image)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// convert id to uuid
	id, err := util.ConvertStringToUUID(req.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the file
	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer file.Close()

	// convert file data to byte
	bytes, err := fileToBytes(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// save file data and name to db
	arg := db.CreateImagesParams{
		ImageName: fileHeader.Filename,
		Data:      bytes,
		Product:   id,
	}

	image, err := server.store.CreateImages(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return a response
	response := newUploadImageResponse(image)

	success(ctx, response)
}

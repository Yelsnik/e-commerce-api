package api

import (
	"net/http"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/token"
	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	Category     string  `json:"category" binding:"required"`
	ProductName  string  `json:"product_name" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Brand        string  `json:"brand"`
	CountInStock int64   `json:"count_in_stock" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
}

type ProductResponse struct {
	Message string     `json:"message"`
	Data    db.Product `json:"data"`
}

func newProductResponse(p db.Product) ProductResponse {
	return ProductResponse{
		Message: "success",
		Data:    p,
	}
}

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateProductsParams{
		Category:     req.Category,
		ProductName:  req.ProductName,
		Description:  req.Description,
		Brand:        util.NewNullString(req.Brand),
		CountInStock: req.CountInStock,
		Price:        req.Price,
		UserID:       authPayload.User_ID,
	}

	product, err := server.store.CreateProducts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newProductResponse(product)

	success(ctx, response)
}

type getProductRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) getProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.GetProducts(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newProductResponse(product)

	success(ctx, response)
}

type listProductRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listProductResponse struct {
	Message string       `json:"message"`
	Data    []db.Product `json:"data"`
}

func newlistProductResponse(p []db.Product) listProductResponse {
	return listProductResponse{
		Message: "success",
		Data:    p,
	}
}

func (server *Server) listProduct(ctx *gin.Context) {
	var req listProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.PageID == 0 && req.PageSize == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "you have not entered a query parameter"})
		return
	}

	arg := db.ListProductsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	products, err := server.store.ListProducts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newlistProductResponse(products)
	success(ctx, response)

}

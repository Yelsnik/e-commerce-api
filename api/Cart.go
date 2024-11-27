package api

import (
	"database/sql"
	"fmt"

	//"fmt"
	"net/http"

	_ "github.com/lib/pq"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/token"
	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/gin-gonic/gin"
)

type addToCartRequest struct {
	ProductID string `uri:"id" binding:"required"`
	Quantity  int64  `json:"quantity" binding:"required,min=1"`
}

type addToCartResponse struct {
	CartItem db.Cartitem `json:"cart_item"`
	Cart     db.Cart     `json:"cart"`
}

func newCartResponse(result db.CartTxResult) addToCartResponse {
	return addToCartResponse{
		CartItem: result.CartItem,
		Cart:     result.Cart,
	}
}

func (server *Server) addToCartApi(ctx *gin.Context) {
	var req addToCartRequest
	err := ctx.ShouldBindUri(&req)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		fmt.Println(req.ProductID, req.Quantity)
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// convert id to uuid
	id, err := util.ConvertStringToUUID(req.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	product, err := server.store.GetProducts(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check if cart items exists
	_, err = server.store.GetCartitemsByProductID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No cart item found; proceed as normal
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		// If no error, a cart item exists
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Product is already added to your cart"})
		return
	}

	// get carts or create one for the user if no cart
	cart, err := server.store.GetCartsByUserID(ctx, authPayload.User_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			arg := db.CreateCartsParams{
				UserID:     authPayload.User_ID,
				TotalPrice: 0,
			}

			// create cart
			cart, err := server.store.CreateCarts(ctx, arg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			if req.Quantity == 0 {
				req.Quantity = 1
			}

			if req.Quantity > product.CountInStock {
				message := fmt.Sprintf("The number of %s you want to buy has exceeded the quantity available for sale", product.ProductName)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": message})
				return
			}

			subTotal := product.Price * float64(1)
			args := db.CreateCartitemsParams{
				Cart:     cart.ID,
				Product:  id,
				Quantity: req.Quantity,
				Price:    product.Price,
				Currency: product.Currency,
				SubTotal: subTotal,
			}

			result, err := server.store.AddToCartTx(ctx, args)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			response := newCartResponse(result)

			success(ctx, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}

	// create cart item
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	if req.Quantity > product.CountInStock {
		message := fmt.Sprintf("The number of %s you want to buy has exceeded the quantity available for sale", product.ProductName)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": message})
		return
	}

	subTotal := product.Price * float64(1)
	arg := db.CreateCartitemsParams{
		Cart:     cart.ID,
		Product:  id,
		Quantity: req.Quantity,
		Price:    product.Price,
		Currency: product.Currency,
		SubTotal: subTotal,
	}

	result, err := server.store.AddToCartTx(ctx, arg)

	response := newCartResponse(result)

	success(ctx, response)
}

type updateCartRequest struct {
	ProductID string `uri:"id" binding:"required"`
	Quantity  int64  `json:"quantity" binding:"required,min=1"`
}

func (server *Server) updateCartApi(ctx *gin.Context) {
	var req updateCartRequest
	err := ctx.ShouldBindUri(&req)
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// convert id to uuid
	id, err := util.ConvertStringToUUID(req.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get cart items by id
	cartItem, err := server.store.GetCartitemsByProductID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// calculate subtotal
	subTotal := float64(req.Quantity) * cartItem.Price
	arg := db.UpdateCartitemsParams{
		ID:       id,
		Quantity: req.Quantity,
		SubTotal: subTotal,
	}

	// update cart items using db transaction
	result, err := server.store.UpdateCartTx(ctx, cartItem.ID, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newCartResponse(result)

	success(ctx, response)
}

type getAllCartsApiResponse struct {
	CartItems []db.Cartitem `json:"cart_items"`
}

func newGetAllCartsApiResponse(cartItems []db.Cartitem) getAllCartsApiResponse {
	return getAllCartsApiResponse{
		CartItems: cartItems,
	}
}

func (server *Server) getAllCartsApi(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	cart, err := server.store.GetCartsByUserID(ctx, authPayload.User_ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cartItems, err := server.store.GetCartitemsByCartID(ctx, cart.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newGetAllCartsApiResponse(cartItems)

	success(ctx, response)
}

type removeCartItemApiRequest struct {
	ProductID string `uri:"id" binding:"required"`
}

func (server *Server) removeCartItemApi(ctx *gin.Context) {
	var req removeCartItemApiRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// convert id to uuid
	id, err := util.ConvertStringToUUID(req.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cartItem, err := server.store.GetCartitemsByProductID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.RemoveCartTx(ctx, cartItem.ID, cartItem.Cart)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success!"})

}

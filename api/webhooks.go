package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/webhook"
)

type paymentSuccessResponse struct {
	Payment db.Payment `json:"payment"`
	Order   db.Order   `json:"order"`
}

func (server *Server) stripeWebhook(ctx *gin.Context) {
	const MaxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	webHookSecret := server.config.WebhookSigningKey
	event, err := webhook.ConstructEvent(payload, ctx.GetHeader("Stripe-Signature"),
		webHookSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		server.handlePaymentIfSuccessful(ctx, &paymentIntent)

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	//ctx.JSON(http.StatusOK, gin.H{"message": "successful"})
}

func (server *Server) handlePaymentIfSuccessful(ctx *gin.Context, paymentIntent *stripe.PaymentIntent) {
	params := &stripe.CustomerParams{}

	customer, err := customer.Get(paymentIntent.Customer.ID, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, customer.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// convert cart item id to uuid
	id, err := uuid.Parse(paymentIntent.Metadata["cartItemID"])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cartItem, err := server.store.GetCartitem(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "unable to find cart item", "error": err})
		return
	}

	// convert string to float
	floatValue, err := strconv.ParseFloat(paymentIntent.Metadata["amount"], 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "unable to convert amount to float", "error": err})
		return
	}

	amount := floatValue / 100

	arg := db.OrderTxParams{
		PaymentIntent:   paymentIntent.ID,
		UserName:        user.Name,
		UserID:          user.ID,
		TotalPrice:      amount,
		DeliveryAddress: paymentIntent.Metadata["deliveryAddress"],
		Country:         paymentIntent.Metadata["country"],
		PaymentStatus:   string(paymentIntent.Status),
		OrderStatus:     "processing",
	}

	result, err := server.store.CreateOrderTx(ctx, cartItem.Cart, cartItem.ID, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "unable to create order", "error": err})
		return
	}

	response := paymentSuccessResponse{
		Payment: result.Payment,
		Order:   result.Order,
	}

	success(ctx, response)
}

type updateOrderStatusRequest struct {
	OrderID string `uri:"order_id"`
	Status  string `json:"status"`
}

func (server *Server) updateOrderStatus(ctx *gin.Context) {
	var req updateOrderStatusRequest
	err := ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	id, err := uuid.Parse(req.OrderID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateOrdersParams{
		ID:     id,
		Status: req.Status,
	}

	order, err := server.store.UpdateOrders(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	success(ctx, order)

}

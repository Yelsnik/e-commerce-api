package api

import (
	"fmt"
	"net/http"
	"strconv"

	db "github.com/Yelsnik/e-commerce-api/db/sqlc"
	"github.com/Yelsnik/e-commerce-api/token"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

type checkoutRequest struct {
	Amount          float64 `json:"amount" binding:"required"`
	Country         string  `json:"country" binding:"required"`
	DeliveryAddress string  `json:"delivery_address" binding:"required"`
	Currency        string  `json:"currency" binding:"required"`
	CartItemID      string  `uri:"id"`
}

type checkoutResponse struct {
	ClientSecret string       `json:"client_secret"`
	Payment      db.Payment   `json:"payment"`
	User         userResponse `json:"user"`
}

func (server *Server) createStripeCustomer(userName, userEmail string) (*stripe.Customer, error) {
	email := fmt.Sprintf("email: '%s'", userEmail)
	searchParams := &stripe.CustomerSearchParams{
		SearchParams: stripe.SearchParams{
			Query: email,
		},
	}

	result := customer.Search(searchParams)
	data := result.CustomerSearchResult().Data

	for i := 0; i < len(data); i++ {
		if data[i].Email == userEmail {
			fmt.Println(data[i])
			return data[i], nil
		}
	}

	params := &stripe.CustomerParams{
		Name:  stripe.String(userName),
		Email: stripe.String(userEmail),
	}

	customer, err := customer.New(params)

	return customer, err
}

func (server *Server) checkout(ctx *gin.Context) {
	var req checkoutRequest
	_ = ctx.ShouldBindJSON(&req)
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUser(ctx, authPayload.User_ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	customer, err := server.createStripeCustomer(user.Name, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	amount := strconv.FormatFloat(req.Amount, 'f', 2, 64)
	metaData := map[string]string{
		"amount":          amount,
		"country":         req.Country,
		"deliveryAddress": req.DeliveryAddress,
		"currency":        req.Currency,
		"cartItemID":      req.CartItemID,
		"userID":          user.ID.String(),
	}

	params := &stripe.PaymentIntentParams{
		Customer: stripe.String(customer.ID),
		Amount:   stripe.Int64(int64(req.Amount)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled:        stripe.Bool(true),
			AllowRedirects: stripe.String("never"),
		},
		Metadata: metaData,
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	amt := req.Amount / 100

	payment, err := server.store.CreatePayments(ctx, db.CreatePaymentsParams{
		ID:       intent.ID,
		Amount:   amt,
		Currency: req.Currency,
		Status:   "processing",
		UserID:   user.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := checkoutResponse{
		ClientSecret: intent.ClientSecret,
		Payment:      payment,
		User:         newSignUpResponse(user),
	}

	success(ctx, response)

}

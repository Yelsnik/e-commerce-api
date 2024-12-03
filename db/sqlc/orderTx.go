package db

import (
	"context"

	"github.com/google/uuid"
)

type OrderTxResult struct {
	Order     Order     `json:"order"`
	OrderItem Orderitem `json:"order_item"`
	Cart      Cart      `json:"cart"`
	CartItem  Cartitem  `json:"cart_item"`
	Payment   Payment   `json:"payment"`
}

type OrderTxParams struct {
	PaymentIntent   string    `json:"payment_intent"`
	UserName        string    `json:"user_name"`
	UserID          uuid.UUID `json:"user_id"`
	TotalPrice      float64   `json:"total_price"`
	DeliveryAddress string    `json:"delivery_address"`
	Country         string    `json:"country"`
	PaymentStatus   string    `json:"status"`
}

func (store *SQLStore) CreateOrderTx(ctx context.Context, cartID, cartItemID uuid.UUID, arg OrderTxParams) (OrderTxResult, error) {
	var result OrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// get the cart item that has been paid for
		result.CartItem, err = q.GetCartitemsForUpdate(ctx, cartItemID)
		if err != nil {
			return err
		}

		// create payment after it has been successful
		result.Payment, err = q.UpdatePaymentStatus(ctx, UpdatePaymentStatusParams{
			ID:     arg.PaymentIntent,
			Status: arg.PaymentStatus,
		})
		if err != nil {
			return err
		}

		// create order
		result.Order, err = q.CreateOrders(ctx, CreateOrdersParams{
			UserName:        arg.UserName,
			UserID:          arg.UserID,
			TotalPrice:      arg.TotalPrice,
			DeliveryAddress: arg.DeliveryAddress,
			Country:         arg.Country,
		})
		if err != nil {
			return err
		}

		// get the product associated with the cart item
		product, err := q.GetProducts(ctx, result.CartItem.Product)
		if err != nil {
			return err
		}

		// add the product to the order items
		result.OrderItem, err = q.CreateOrderitems(ctx, CreateOrderitemsParams{
			ItemName:     product.ProductName,
			ItemSubTotal: result.CartItem.SubTotal,
			Quantity:     result.CartItem.Quantity,
			ItemID:       product.ID,
			OrderID:      result.Order.ID,
		})
		if err != nil {
			return err
		}

		// remove the cart item from the cart
		err = q.DeleteCartitems(ctx, result.CartItem.ID)
		if err != nil {
			return err
		}

		// get cart for update
		result.Cart, err = q.GetCartsForUpdate(ctx, cartID)
		if err != nil {
			return err
		}

		// add the prices of the items remaining in the cart
		total, err := q.AddSubtotalPrice(ctx, cartID)
		if err != nil {
			return err
		}

		// update the cart
		result.Cart, err = q.UpdateCarts(ctx, UpdateCartsParams{
			ID:         result.Cart.ID,
			TotalPrice: total,
		})
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
}

type Cartitem struct {
	ID       uuid.UUID `json:"id"`
	Cart     uuid.UUID `json:"cart"`
	Product  uuid.UUID `json:"product"`
	Quantity int64     `json:"quantity"`
	// must be positive
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	SubTotal  float64   `json:"sub_total"`
	CreatedAt time.Time `json:"created_at"`
}

type Image struct {
	ID        int64     `json:"id"`
	ImageName string    `json:"image_name"`
	Data      []byte    `json:"data"`
	Product   uuid.UUID `json:"product"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID              uuid.UUID `json:"id"`
	UserName        string    `json:"user_name"`
	UserID          uuid.UUID `json:"user_id"`
	TotalPrice      float64   `json:"total_price"`
	DeliveryAddress string    `json:"delivery_address"`
	Country         string    `json:"country"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type Orderitem struct {
	ID           uuid.UUID `json:"id"`
	ItemName     string    `json:"item_name"`
	ItemSubTotal float64   `json:"item_sub_total"`
	Quantity     int64     `json:"quantity"`
	ItemID       uuid.UUID `json:"item_id"`
	OrderID      uuid.UUID `json:"order_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type PasswordResetToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Payment struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID           uuid.UUID      `json:"id"`
	Category     string         `json:"category"`
	ProductName  string         `json:"product_name"`
	Description  string         `json:"description"`
	Brand        sql.NullString `json:"brand"`
	CountInStock int64          `json:"count_in_stock"`
	// must be positive
	Price      float64       `json:"price"`
	Currency   string        `json:"currency"`
	Rating     sql.NullInt64 `json:"rating"`
	IsFeatured sql.NullBool  `json:"is_featured"`
	UserID     uuid.UUID     `json:"user_id"`
	CreatedAt  time.Time     `json:"created_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

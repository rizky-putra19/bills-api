package dto

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type FullfilmentResponse struct {
	OrderId     int    `json:"order_id"`
	OrderStatus string `json:"order_status"`
}

type OrderPayload struct {
	CustomerID   int    `json:"customer_id"`
	ProductID    int    `json:"product_id"`
	ClientNumber string `json:"client_number"`
}

type PaymentMetadata struct {
	OrderID           int     `json:"order_id"`
	ReferenceNumber   string  `json:"reference_number"`
	AccountNumber     string  `json:"account_number"`
	Amount            float64 `json:"amount"`
	LFIToken          string  `json:"lfi_token"` // need in payment for redirect payment web
	PaymentExpiryTime string  `json:"payment_expiry_time"`
}

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

type OrderDetails struct {
	ID                int        `db:"id" json:"id"`
	CustomerID        int        `db:"customer_id" json:"customer_id"`
	OrderStatus       string     `db:"order_status" json:"order_status"`
	PaymentStatus     string     `db:"payment_status" json:"payment_status"`
	OrderCreatedAt    *time.Time `db:"order_created_at" json:"order_created_at"`
	TotalPrice        float64    `db:"total_price" json:"total_price"`
	ClientNumber      string     `db:"client_number" json:"client_number"`
	ProductName       string     `db:"product_name" json:"product_name"`
	FulfillmentStatus string     `db:"fulfillment_status" json:"fulfillment_status"`
	AccountNumber     string     `db:"account_number" json:"account_number"`
	PaymentMethod     string     `db:"payment_method" json:"payment_method"`
	PaymentExpiryTime *time.Time `db:"payment_expiry_time" json:"payment_expiry_time"`
}

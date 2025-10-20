package database

import (
	"time"
)

// Order local order model
type Order struct {
	ID                int       `json:"id"`
	OrderNumber       string    `json:"order_number"`
	CustomerName      string    `json:"customer_name"`
	CustomerEmail     string    `json:"customer_email"`
	ProductID         string    `json:"product_id"`
	ProductTokenID    string    `json:"product_token_id"`
	Quantity          int       `json:"quantity"`
	TotalAmount       string    `json:"total_amount"`
	Status            string    `json:"status"` // pending, paid, cancelled, expired
	ReddioPaymentID   string    `json:"reddio_payment_id"`
	ReddioPayLink     string    `json:"reddio_pay_link"`
	ReddioStatus      string    `json:"reddio_status"`
	TransactionHash   string    `json:"transaction_hash"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PaidAt            *time.Time `json:"paid_at,omitempty"`
}

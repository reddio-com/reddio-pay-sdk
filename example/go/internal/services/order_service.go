package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"order-system/internal/config"
	"order-system/internal/database"
	"time"

	"github.com/reddio-com/reddio-pay-sdk/go-sdk/client"
)

type OrderService struct {
	db     *sql.DB
	config *config.Config
	client *client.Client
}

// NewReddioPayClient creates Reddio Pay SDK client
func NewReddioPayClient(cfg *config.Config) (*client.Client, error) {
	ctx := context.Background()
	return client.NewSDKClient(ctx, cfg.ReddioURL, cfg.ReddioAPIKey)
}

func NewOrderService(db *sql.DB, cfg *config.Config) *OrderService {
	// 初始化 Reddio Pay SDK 客户端
	ctx := context.Background()
	reddioClient, err := client.NewSDKClient(ctx, cfg.ReddioURL, cfg.ReddioAPIKey)
	if err != nil {
		log.Fatal("Failed to initialize Reddio Pay client:", err)
	}

	return &OrderService{
		db:     db,
		config: cfg,
		client: reddioClient,
	}
}

// CreateOrder creates an order (simplified version, mainly for SDK demonstration)
func (s *OrderService) CreateOrder(customerName, customerEmail, productID, productTokenID string, quantity int) (*database.Order, error) {
	// Generate order number
	orderNumber := s.generateOrderNumber()

	// Simplified total amount calculation (in real application, should get accurate price from SDK)
	totalAmount := fmt.Sprintf("%.2f", float64(quantity)*100.0) // Assume $100 per item

	// Create local order
	query := `INSERT INTO orders (order_number, customer_name, customer_email, product_id, product_token_id, quantity, total_amount, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	result, err := s.db.Exec(query, orderNumber, customerName, customerEmail, productID, productTokenID, quantity, totalAmount, "pending", now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get order ID: %w", err)
	}

	order := &database.Order{
		ID:             int(orderID),
		OrderNumber:    orderNumber,
		CustomerName:   customerName,
		CustomerEmail:  customerEmail,
		ProductID:      productID,
		ProductTokenID: productTokenID,
		Quantity:       quantity,
		TotalAmount:    totalAmount,
		Status:         "pending",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Call Reddio Pay SDK to create payment
	log.Printf("Calling Reddio Pay SDK to create external payment for order %s", orderNumber)
	reddioReq := &client.ExternalCreatePaymentRequest{
		ProductID:      productID,
		ProductTokenID: productTokenID,
		Count:          quantity,
	}

	reddioResp, err := s.client.ExternalCreatePayment(reddioReq)
	if err != nil {
		// If Reddio Pay creation fails, update order status
		s.db.Exec("UPDATE orders SET status = ? WHERE id = ?", "failed", orderID)
		return nil, fmt.Errorf("failed to create Reddio Pay payment: %w", err)
	}

	// Update order information
	updateQuery := `UPDATE orders SET reddio_payment_id = ?, reddio_pay_link = ?, reddio_status = ?, updated_at = ? WHERE id = ?`
	_, err = s.db.Exec(updateQuery, reddioResp.PaymentID, reddioResp.PayLink, "created", time.Now(), orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to update order with Reddio Pay info: %w", err)
	}

	order.ReddioPaymentID = reddioResp.PaymentID
	order.ReddioPayLink = reddioResp.PayLink
	order.ReddioStatus = "created"

	log.Printf("Order %s created successfully with Reddio Pay ID: %s", orderNumber, reddioResp.PaymentID)
	return order, nil
}

// GetOrder gets order information
func (s *OrderService) GetOrder(orderID int) (*database.Order, error) {
	query := `SELECT id, order_number, customer_name, customer_email, product_id, product_token_id, quantity, 
			  total_amount, status, reddio_payment_id, reddio_pay_link, reddio_status, transaction_hash, 
			  created_at, updated_at, paid_at FROM orders WHERE id = ?`
	
	row := s.db.QueryRow(query, orderID)
	
	var order database.Order
	var paidAt sql.NullTime
	
	err := row.Scan(&order.ID, &order.OrderNumber, &order.CustomerName, &order.CustomerEmail, 
		&order.ProductID, &order.ProductTokenID, &order.Quantity, &order.TotalAmount, 
		&order.Status, &order.ReddioPaymentID, &order.ReddioPayLink, &order.ReddioStatus, 
		&order.TransactionHash, &order.CreatedAt, &order.UpdatedAt, &paidAt)
	
	if err != nil {
		return nil, err
	}
	
	if paidAt.Valid {
		order.PaidAt = &paidAt.Time
	}
	
	return &order, nil
}

// GetOrderByNumber gets order by order number
func (s *OrderService) GetOrderByNumber(orderNumber string) (*database.Order, error) {
	query := `SELECT id, order_number, customer_name, customer_email, product_id, product_token_id, quantity, 
			  total_amount, status, reddio_payment_id, reddio_pay_link, reddio_status, transaction_hash, 
			  created_at, updated_at, paid_at FROM orders WHERE order_number = ?`
	
	row := s.db.QueryRow(query, orderNumber)
	
	var order database.Order
	var paidAt sql.NullTime
	
	err := row.Scan(&order.ID, &order.OrderNumber, &order.CustomerName, &order.CustomerEmail, 
		&order.ProductID, &order.ProductTokenID, &order.Quantity, &order.TotalAmount, 
		&order.Status, &order.ReddioPaymentID, &order.ReddioPayLink, &order.ReddioStatus, 
		&order.TransactionHash, &order.CreatedAt, &order.UpdatedAt, &paidAt)
	
	if err != nil {
		return nil, err
	}
	
	if paidAt.Valid {
		order.PaidAt = &paidAt.Time
	}
	
	return &order, nil
}

// ListOrders gets order list
func (s *OrderService) ListOrders(page, limit int, status string) ([]*database.Order, int64, error) {
	var orders []*database.Order
	var total int64

	// Build query conditions
	whereClause := ""
	args := []interface{}{}
	if status != "" {
		whereClause = "WHERE status = ?"
		args = append(args, status)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM orders " + whereClause
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Paginated query
	offset := (page - 1) * limit
	query := `SELECT id, order_number, customer_name, customer_email, product_id, product_token_id, quantity, 
			  total_amount, status, reddio_payment_id, reddio_pay_link, reddio_status, transaction_hash, 
			  created_at, updated_at, paid_at FROM orders ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	args = append(args, limit, offset)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var order database.Order
		var paidAt sql.NullTime
		
		err := rows.Scan(&order.ID, &order.OrderNumber, &order.CustomerName, &order.CustomerEmail, 
			&order.ProductID, &order.ProductTokenID, &order.Quantity, &order.TotalAmount, 
			&order.Status, &order.ReddioPaymentID, &order.ReddioPayLink, &order.ReddioStatus, 
			&order.TransactionHash, &order.CreatedAt, &order.UpdatedAt, &paidAt)
		
		if err != nil {
			return nil, 0, err
		}
		
		if paidAt.Valid {
			order.PaidAt = &paidAt.Time
		}
		
		orders = append(orders, &order)
	}

	return orders, total, nil
}

// CheckPaymentStatus checks payment status
func (s *OrderService) CheckPaymentStatus(orderID int) error {
	order, err := s.GetOrder(orderID)
	if err != nil {
		return err
	}

	if order.ReddioPaymentID == "" {
		return fmt.Errorf("no Reddio Pay payment ID found for order %d", orderID)
	}

	// Call Reddio Pay SDK to query payment status
	log.Printf("Calling Reddio Pay SDK to check payment status for payment ID: %s", order.ReddioPaymentID)
	payment, err := s.client.GetPaymentByID(order.ReddioPaymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment status from Reddio Pay: %w", err)
	}

	// Update order status
	now := time.Now()
	if payment.Status == "paid" && order.Status != "paid" {
		query := `UPDATE orders SET status = ?, reddio_status = ?, transaction_hash = ?, paid_at = ?, updated_at = ? WHERE id = ?`
		_, err = s.db.Exec(query, "paid", payment.Status, payment.TransactionHash, now, now, orderID)
	} else {
		query := `UPDATE orders SET reddio_status = ?, updated_at = ? WHERE id = ?`
		_, err = s.db.Exec(query, payment.Status, now, orderID)
	}

	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}


// generateOrderNumber generates order number
func (s *OrderService) generateOrderNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("ORD%08d", rand.Intn(99999999))
}

// parseFloat parses string to float64
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
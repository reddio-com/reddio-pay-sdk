package services

import (
	"fmt"
	"log"
	"time"

	"github.com/reddio-com/reddio-pay-sdk/go-sdk/client"
)

// SDKDemoService SDK demonstration service
type SDKDemoService struct {
	client *client.Client
}

// NewSDKDemoService creates SDK demonstration service
func NewSDKDemoService(client *client.Client) *SDKDemoService {
	return &SDKDemoService{
		client: client,
	}
}

// DemoAllAPIs demonstrates all SDK APIs
func (s *SDKDemoService) DemoAllAPIs() error {
	log.Println("=== Reddio Pay Go SDK API Demonstration Started ===")
	
	// 1. Demonstrate getting product list
	if err := s.demoListProducts(); err != nil {
		log.Printf("Failed to demonstrate list products: %v", err)
		log.Println("Continuing with other APIs...")
	}
	
	// 2. Demonstrate creating product
	productID, err := s.demoCreateProduct()
	if err != nil {
		log.Printf("Failed to demonstrate create product: %v", err)
		log.Println("Continuing with other APIs...")
		productID = "demo-product-id" // Use mock ID
	}
	
	// 3. Demonstrate getting product info
	if err := s.demoActivateProduct(productID); err != nil {
		log.Printf("Failed to demonstrate get product info: %v", err)
		log.Println("Continuing with other APIs...")
	}
	
	// 4. Demonstrate adding product token
	productTokenID, err := s.demoCreateProductToken(productID)
	if err != nil {
		log.Printf("Failed to demonstrate add product token: %v", err)
		log.Println("Continuing with other APIs...")
		productTokenID = "demo-product-token-id" // Use mock ID
	}
	
	// 5. Demonstrate creating external payment
	paymentID, err := s.demoExternalCreatePayment(productID, productTokenID)
	if err != nil {
		log.Printf("Failed to demonstrate create external payment: %v", err)
		log.Println("Continuing with other APIs...")
		paymentID = "demo-payment-id" // Use mock ID
	}
	
	// 6. Demonstrate querying payment status
	if err := s.demoGetPaymentByID(paymentID); err != nil {
		log.Printf("Failed to demonstrate get payment by ID: %v", err)
		log.Println("Continuing with other APIs...")
	}
	
	log.Println("=== Reddio Pay Go SDK API Demonstration Completed ===")
	return nil
}

// demoListProducts demonstrates getting product list
func (s *SDKDemoService) demoListProducts() error {
	log.Println("\n--- 1. Demonstrate List Products ---")
	
	productsResp, err := s.client.ListProducts()
	if err != nil {
		return fmt.Errorf("failed to get product list: %w", err)
	}
	
	log.Printf("Successfully retrieved %d products:", len(productsResp.Products))
	for i, product := range productsResp.Products {
		log.Printf("  Product %d: ID=%s, Name=%s, Active=%t", 
			i+1, product.ProductID, product.Name, product.Active)
	}
	
	return nil
}

// demoCreateProduct demonstrates creating product
func (s *SDKDemoService) demoCreateProduct() (string, error) {
	log.Println("\n--- 2. Demonstrate Create Product ---")
	
	createReq := &client.CreateProductRequest{
		Name:             "SDK Demo Product",
		Description:      "This is a demo product created through Go SDK",
		Content:          "Product content description",
		TokenIDList:      []string{"0x1234567890abcdef1234567890abcdef12345678"},
		Price:            "1000000000000000000", // 1 ETH in wei
		RecipientAddress: "0xabcdef1234567890abcdef1234567890abcdef12",
	}
	
	createResp, err := s.client.CreateProduct(createReq)
	if err != nil {
		return "", fmt.Errorf("failed to create product: %w", err)
	}
	
	log.Printf("Successfully created product: ID=%s, Name=%s", createResp.Product.ProductID, createResp.Product.Name)
	return createResp.Product.ProductID, nil
}

// demoActivateProduct demonstrates getting product info (shows product status)
func (s *SDKDemoService) demoActivateProduct(productID string) error {
	log.Println("\n--- 3. Demonstrate Get Product Info ---")
	
	product, err := s.client.GetProduct(productID)
	if err != nil {
		return fmt.Errorf("failed to get product info: %w", err)
	}
	
	log.Printf("Product info: ID=%s, Name=%s, Active=%t", product.ProductID, product.Name, product.Active)
	return nil
}

// demoCreateProductToken demonstrates adding product token
func (s *SDKDemoService) demoCreateProductToken(productID string) (string, error) {
	log.Println("\n--- 4. Demonstrate Add Product Token ---")
	
	addTokenReq := &client.AddProductTokenRequest{
		TokenID:          "0x9876543210fedcba9876543210fedcba98765432", // Example token ID
		Price:            "2000000000000000000", // 2 ETH in wei
		RecipientAddress: "0xfedcba9876543210fedcba9876543210fedcba98", // Example recipient address
	}
	
	addTokenResp, err := s.client.AddProductToken(productID, addTokenReq)
	if err != nil {
		return "", fmt.Errorf("failed to add product token: %w", err)
	}
	
	log.Printf("Successfully added product token: ID=%s, TokenID=%s, Price=%s", 
		addTokenResp.ProductToken.ProductTokenID, addTokenResp.ProductToken.TokenID, addTokenResp.ProductToken.Price)
	return addTokenResp.ProductToken.ProductTokenID, nil
}

// demoExternalCreatePayment demonstrates creating external payment
func (s *SDKDemoService) demoExternalCreatePayment(productID, productTokenID string) (string, error) {
	log.Println("\n--- 5. Demonstrate Create External Payment ---")
	
	createPaymentReq := &client.ExternalCreatePaymentRequest{
		ProductID:      productID,
		ProductTokenID: productTokenID,
		Count:          1,
	}
	
	createPaymentResp, err := s.client.ExternalCreatePayment(createPaymentReq)
	if err != nil {
		return "", fmt.Errorf("failed to create external payment: %w", err)
	}
	
	log.Printf("Successfully created external payment: PaymentID=%s, PayLink=%s", 
		createPaymentResp.PaymentID, createPaymentResp.PayLink)
	return createPaymentResp.PaymentID, nil
}

// demoGetPaymentByID demonstrates querying payment status
func (s *SDKDemoService) demoGetPaymentByID(paymentID string) error {
	log.Println("\n--- 6. Demonstrate Get Payment By ID ---")
	
	payment, err := s.client.GetPaymentByID(paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment status: %w", err)
	}
	
	log.Printf("Payment status: ID=%s, Status=%s, TransactionHash=%s", 
		payment.PaymentID, payment.Status, payment.TransactionHash)
	return nil
}

// DemoWithDelay demonstrates with delays (for showing async operations)
func (s *SDKDemoService) DemoWithDelay() error {
	log.Println("\n=== SDK API Demonstration with Delays ===")
	
	// Create product
	productID, err := s.demoCreateProduct()
	if err != nil {
		return err
	}
	
	// Wait for a while
	log.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	
	// Get product info
	if err := s.demoActivateProduct(productID); err != nil {
		return err
	}
	
	// Wait for a while
	log.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	
	// Add product token
	productTokenID, err := s.demoCreateProductToken(productID)
	if err != nil {
		return err
	}
	
	// Wait for a while
	log.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	
	// Create payment
	paymentID, err := s.demoExternalCreatePayment(productID, productTokenID)
	if err != nil {
		return err
	}
	
	// Wait for a while
	log.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	
	// Query payment status
	if err := s.demoGetPaymentByID(paymentID); err != nil {
		return err
	}
	
	log.Println("=== SDK API Demonstration with Delays Completed ===")
	return nil
}
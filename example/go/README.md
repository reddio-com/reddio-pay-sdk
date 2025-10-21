# Reddio Pay Go SDK Example

This example demonstrates how to use the `github.com/reddio-com/reddio-pay-sdk/go-sdk` SDK in a Go application. The example shows how to integrate Reddio Pay SDK into a backend service for order management and payment processing.

## Overview

This example server showcases the integration of Reddio Pay Go SDK with the following features:

- **SDK Initialization**: How to initialize the Reddio Pay SDK client
- **Token Management**: Getting currently supported token types using ListTokens API
- **Product Management**: Creating and managing products using SDK APIs
- **Payment Processing**: Creating external payments and checking payment status
- **Order Management**: A simple order system that demonstrates SDK usage

## SDK APIs Demonstrated

The example demonstrates the following Reddio Pay SDK APIs:

### 1. SDK Client Initialization

```go
import "github.com/reddio-com/reddio-pay-sdk/go-sdk/client"

// Initialize SDK client
ctx := context.Background()
client, err := client.NewSDKClient(ctx, "https://reddio-service-prod.reddio.com", apiKey)
if err != nil {
    log.Fatal("Failed to initialize Reddio Pay client:", err)
}
```

### 2. List Tokens

```go
// Get all currently supported token types
tokensResp, err := client.ListTokens()
if err != nil {
    return fmt.Errorf("failed to get token list: %w", err)
}

log.Printf("Successfully retrieved %d supported tokens:", tokensResp.Count)
for i, token := range tokensResp.Tokens {
    log.Printf("  Token %d: ID=%s, Name=%s, Symbol=%s, Chain=%s, Active=%t", 
        i+1, token.TokenID, token.Name, token.Symbol, token.ChainName, token.IsActive)
}
```

**Purpose**: This API returns all currently supported token types that can be used for payments. It shows available tokens with their details including token ID, name, symbol, chain information, and active status.

### 3. List Products

```go
// Get all products
productsResp, err := client.ListProducts()
if err != nil {
    return fmt.Errorf("failed to get product list: %w", err)
}

log.Printf("Successfully retrieved %d products:", len(productsResp.Products))
for i, product := range productsResp.Products {
    log.Printf("  Product %d: ID=%s, Name=%s, Active=%t", 
        i+1, product.ProductID, product.Name, product.Active)
}
```

### 4. Create Product

```go
// Create a new product
createReq := &client.CreateProductRequest{
    Name:             "Demo Product",
    Description:      "This is a demo product created through Go SDK",
    Content:          "Product content description",
    TokenIDList:      []string{"0x1234567890abcdef1234567890abcdef12345678"},
    Price:            "1000000000000000000", // 1 ETH in wei
    RecipientAddress: "0xabcdef1234567890abcdef1234567890abcdef12",
}

createResp, err := client.CreateProduct(createReq)
if err != nil {
    return fmt.Errorf("failed to create product: %w", err)
}

log.Printf("Successfully created product: ID=%s, Name=%s", 
    createResp.Product.ProductID, createResp.Product.Name)
```

### 5. Get Product Information

```go
// Get product details
product, err := client.GetProduct(productID)
if err != nil {
    return fmt.Errorf("failed to get product info: %w", err)
}

log.Printf("Product info: ID=%s, Name=%s, Active=%t", 
    product.ProductID, product.Name, product.Active)
```

### 6. Add Product Token

```go
// Add a token to an existing product
addTokenReq := &client.AddProductTokenRequest{
    TokenID:          "0x9876543210fedcba9876543210fedcba98765432",
    Price:            "2000000000000000000", // 2 ETH in wei
    RecipientAddress: "0xfedcba9876543210fedcba9876543210fedcba98",
}

addTokenResp, err := client.AddProductToken(productID, addTokenReq)
if err != nil {
    return fmt.Errorf("failed to add product token: %w", err)
}

log.Printf("Successfully added product token: ID=%s, TokenID=%s, Price=%s", 
    addTokenResp.ProductToken.ProductTokenID, 
    addTokenResp.ProductToken.TokenID, 
    addTokenResp.ProductToken.Price)
```

### 7. Create External Payment

```go
// Create an external payment for an order
createPaymentReq := &client.ExternalCreatePaymentRequest{
    ProductID:      productID,
    ProductTokenID: productTokenID,
    Count:          1,
}

createPaymentResp, err := client.ExternalCreatePayment(createPaymentReq)
if err != nil {
    return fmt.Errorf("failed to create external payment: %w", err)
}

log.Printf("Successfully created external payment: PaymentID=%s, PayLink=%s", 
    createPaymentResp.PaymentID, createPaymentResp.PayLink)
```

### 8. Get Payment Status

```go
// Query payment status
payment, err := client.GetPaymentByID(paymentID)
if err != nil {
    return fmt.Errorf("failed to get payment status: %w", err)
}

log.Printf("Payment status: ID=%s, Status=%s, TransactionHash=%s", 
    payment.PaymentID, payment.Status, payment.TransactionHash)
```

## Project Structure

```
example/go/
├── main.go                   # Main entry point with SDK initialization
├── go.mod                    # Go module dependencies
└── internal/
    ├── config/
    │   └── config.go         # Configuration management
    ├── database/
    │   ├── database.go       # Database setup
    │   └── models.go         # Data models
    ├── handlers/
    │   └── order_handler.go  # HTTP handlers
    └── services/
        ├── order_service.go  # Business logic with SDK integration
        └── sdk_demo.go       # SDK API demonstration service
```

## Key Integration Points

### 1. SDK Client in Order Service

The `OrderService` integrates the Reddio Pay SDK to create payments:

```go
type OrderService struct {
    db     *sql.DB
    client *client.Client
}

func (s *OrderService) CreateOrder(customerName, customerEmail, productID, productTokenID string, quantity int) (*database.Order, error) {
    // ... create local order ...
    
    // Call Reddio Pay SDK to create payment
    reddioReq := &client.ExternalCreatePaymentRequest{
        ProductID:      productID,
        ProductTokenID: productTokenID,
        Count:          quantity,
    }

    reddioResp, err := s.client.ExternalCreatePayment(reddioReq)
    if err != nil {
        return nil, fmt.Errorf("failed to create Reddio Pay payment: %w", err)
    }

    // Update order with payment information
    order.ReddioPaymentID = reddioResp.PaymentID
    order.ReddioPayLink = reddioResp.PayLink
    order.ReddioStatus = "created"
    
    return order, nil
}
```

### 2. SDK Demonstration Service

The `SDKDemoService` showcases all available SDK APIs:

```go
type SDKDemoService struct {
    client *client.Client
}

func (s *SDKDemoService) DemoAllAPIs() error {
    // Demonstrates all SDK APIs in sequence
    // 1. List Products
    // 2. Create Product
    // 3. Get Product Info
    // 4. Add Product Token
    // 5. Create External Payment
    // 6. Get Payment Status
}
```

### 3. Configuration

The SDK is configured with API key and URL:

```go
type Config struct {
    ReddioAPIKey string
    ReddioURL    string
}

func Load() *Config {
    return &Config{
        ReddioAPIKey: getEnv("REDDIO_API_KEY", "your_api_key"),
        ReddioURL:    getEnv("REDDIO_URL", "https://reddio-service-prod.reddio.com"),
    }
}
```

## Running the Example

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Set Environment Variables (Optional)

```bash
export REDDIO_API_KEY="your_api_key"
export REDDIO_URL="https://reddio-service-prod.reddio.com"
```

### 3. Run the Example

```bash
go run main.go
```

The example will:
1. Initialize the Reddio Pay SDK client
2. Demonstrate all SDK APIs with detailed logging
3. Start a REST API server for order management
4. Show how to integrate SDK calls in business logic

## SDK Integration Benefits

This example demonstrates how the Reddio Pay SDK enables:

- **Easy Integration**: Simple client initialization and API calls
- **Payment Processing**: Seamless payment creation and status tracking
- **Product Management**: Complete product lifecycle management
- **Error Handling**: Proper error handling and logging
- **Flexibility**: Easy to integrate into existing Go applications

## Dependencies

- `github.com/reddio-com/reddio-pay-sdk/go-sdk` - Reddio Pay Go SDK
- `github.com/gorilla/mux` - HTTP router
- `database/sql` - Database operations
- `github.com/mattn/go-sqlite3` - SQLite driver

## Error Handling

The example shows proper error handling patterns:

```go
// SDK API calls with error handling
resp, err := client.ExternalCreatePayment(req)
if err != nil {
    // Handle SDK errors appropriately
    return fmt.Errorf("failed to create payment: %w", err)
}

// Log successful operations
log.Printf("Successfully created payment: %s", resp.PaymentID)
```

This example serves as a comprehensive guide for integrating the Reddio Pay Go SDK into your Go applications.
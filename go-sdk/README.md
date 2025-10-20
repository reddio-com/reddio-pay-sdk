# Reddio Pay Go SDK

Reddio Pay Go SDK is a Go client library for interacting with Reddio Pay services. It provides a complete API interface supporting account management, token management, product management, payment processing, and more.

## Installation

```bash
go get github.com/reddio-com/reddio-pay-sdk/go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/reddio-com/reddio-pay-sdk/go-sdk/client"
)

func main() {
    // Create client
    ctx := context.Background()
    client, err := client.NewSDKClient(ctx, "https://reddio-service-prod.reddio.com", "your-api-key")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Use client for API calls
    // ...
}
```

## Client Initialization

### NewSDKClient

Creates a new SDK client instance.

```go
func NewSDKClient(ctx context.Context, url string, apiKey string) (*Client, error)
```

**Parameters:**
- `ctx`: Context object
- `url`: API server address
- `apiKey`: API key

**Returns:**
- `*Client`: Client instance
- `error`: Error information

**Example:**
```go
ctx := context.Background()
client, err := client.NewSDKClient(ctx, "https://reddio-service-prod.reddio.com", "your-api-key")
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## API Reference

### Account Management

#### GetAccountInfo

Gets current account information.

```go
func (c *Client) GetAccountInfo() (*AccountResponse, error)
```

**Returns:**
- `*AccountResponse`: Account information
- `error`: Error information

**Example:**
```go
account, err := client.GetAccountInfo()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Account Email: %s\n", account.Email)
fmt.Printf("Company Name: %s\n", account.CompanyName)
```

#### UpdateWebhook

Updates the account's webhook address.

```go
func (c *Client) UpdateWebhook(req *UpdateWebhookRequest) (*UpdateWebhookResponse, error)
```

**Parameters:**
- `req`: Request object containing new webhook address

**Example:**
```go
req := &client.UpdateWebhookRequest{
    Webhook: "https://your-domain.com/webhook",
}
response, err := client.UpdateWebhook(req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Update Result: %s\n", response.Message)
```

#### ListAccountAddresses

Gets all address information for the account.

```go
func (c *Client) ListAccountAddresses() ([]*AccountAddress, error)
```

**Returns:**
- `[]*AccountAddress`: Address list
- `error`: Error information

**Example:**
```go
addresses, err := client.ListAccountAddresses()
if err != nil {
    log.Fatal(err)
}
for _, addr := range addresses {
    fmt.Printf("Address: %s, Token: %s\n", addr.RecipientAddress, addr.TokenID)
}
```

### Token Management

#### ListTokens

Gets all available tokens list.

```go
func (c *Client) ListTokens() (*ListTokensResponse, error)
```

**Returns:**
- `*ListTokensResponse`: Token list response
- `error`: Error information

**Example:**
```go
tokens, err := client.ListTokens()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total tokens: %d\n", tokens.Count)
for _, token := range tokens.Tokens {
    fmt.Printf("Token: %s (%s)\n", token.Name, token.Symbol)
}
```

### Product Management

#### CreateProduct

Creates a new product.

```go
func (c *Client) CreateProduct(req *CreateProductRequest) (*CreateProductResponse, error)
```

**Parameters:**
- `req`: Product creation request object

**Example:**
```go
req := &client.CreateProductRequest{
    Name:             "My Product",
    Description:      "Product description",
    Content:          "Product content",
    TokenIDList:      []string{"token1", "token2"},
    Price:            "100.0",
    RecipientAddress: "0x123...",
}
response, err := client.CreateProduct(req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Product created successfully: %s\n", response.Product.ProductID)
```

#### ListProducts

Gets product list.

```go
func (c *Client) ListProducts() (*ListProductsResponse, error)
```

**Example:**
```go
products, err := client.ListProducts()
if err != nil {
    log.Fatal(err)
}
for _, product := range products.Products {
    fmt.Printf("Product: %s\n", product.Name)
}
```

#### GetProduct

Gets product details by product ID.

```go
func (c *Client) GetProduct(productID string) (*Product, error)
```

**Parameters:**
- `productID`: Product ID

**Example:**
```go
product, err := client.GetProduct("product123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Product Name: %s\n", product.Name)
```

#### AddProductToken

Adds a token to a product.

```go
func (c *Client) AddProductToken(productID string, req *AddProductTokenRequest) (*AddProductTokenResponse, error)
```

**Parameters:**
- `productID`: Product ID
- `req`: Add token request object

**Example:**
```go
req := &client.AddProductTokenRequest{
    TokenID:          "token123",
    Price:            "50.0",
    RecipientAddress: "0x456...",
}
response, err := client.AddProductToken("product123", req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Token added successfully: %s\n", response.ProductToken.ProductTokenID)
```

#### GetProductTokenStatus

Gets product token status.

```go
func (c *Client) GetProductTokenStatus(productID string) (*GetProductTokenStatusResponse, error)
```

**Parameters:**
- `productID`: Product ID

**Example:**
```go
status, err := client.GetProductTokenStatus("product123")
if err != nil {
    log.Fatal(err)
}
for _, s := range status.Status {
    fmt.Printf("Token: %s, Sales Count: %d\n", s.TokenName, s.TotalSaleCount)
}
```

### Payment Management

#### ListPaymentsByAccountAndProductID

Gets payment list by product ID.

```go
func (c *Client) ListPaymentsByAccountAndProductID(productID string) (*ListPaymentsResponse, error)
```

**Parameters:**
- `productID`: Product ID

**Example:**
```go
payments, err := client.ListPaymentsByAccountAndProductID("product123")
if err != nil {
    log.Fatal(err)
}
for _, payment := range payments.Payments {
    fmt.Printf("Payment ID: %s, Status: %s\n", payment.PaymentID, payment.Status)
}
```

#### ListPaymentsByAccount

Gets all payment records for the account (with pagination support).

```go
func (c *Client) ListPaymentsByAccount(limit, offset int) (*ListPaymentsResponseWithPagination, error)
```

**Parameters:**
- `limit`: Number of records per page
- `offset`: Offset

**Example:**
```go
// Get first page with 10 records per page
payments, err := client.ListPaymentsByAccount(10, 0)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Total records: %d\n", payments.TotalCount)
fmt.Printf("Current page records: %d\n", len(payments.Payments))
```

#### GetPaymentByID

Gets payment details by payment ID.

```go
func (c *Client) GetPaymentByID(paymentID string) (*Payment, error)
```

**Parameters:**
- `paymentID`: Payment ID

**Example:**
```go
payment, err := client.GetPaymentByID("payment123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Payment Status: %s\n", payment.Status)
fmt.Printf("Total Amount: %s\n", payment.TotalAmount)
```

#### ExternalSendNotifyForPaymentSuccess

Sets up payment success notification.

```go
func (c *Client) ExternalSendNotifyForPaymentSuccess(req *ExternalSendNotifyForPaymentSuccessRequest) (*SendNotifyForPaymentSuccessResponse, error)
```

**Parameters:**
- `req`: Notification setup request object

**Example:**
```go
req := &client.ExternalSendNotifyForPaymentSuccessRequest{
    PaymentID: "payment123",
    Email:     "user@example.com",
}
response, err := client.ExternalSendNotifyForPaymentSuccess(req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Notification setup result: %s\n", response.Message)
```

#### ExternalCreatePayment

Creates a new external payment.

```go
func (c *Client) ExternalCreatePayment(req *ExternalCreatePaymentRequest) (*ExternalCreatePaymentResponse, error)
```

**Parameters:**
- `req`: External payment creation request object

**Example:**
```go
req := &client.ExternalCreatePaymentRequest{
    ProductID:      "product123",
    ProductTokenID: "token456",
    Count:          1,
}
response, err := client.ExternalCreatePayment(req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Payment created successfully: %s\n", response.PaymentID)
fmt.Printf("Payment link: %s\n", response.PayLink)
fmt.Printf("Contract address: %s\n", response.ContractAddress)
```

## Data Structures

### Account Related

#### AccountResponse
```go
type AccountResponse struct {
    Email       string `json:"email"`
    Webhook     string `json:"webhook,omitempty"`
    CompanyName string `json:"company_name,omitempty"`
    CompanyURL  string `json:"company_url,omitempty"`
    Activated   bool   `json:"activated"`
    CreatedAt   string `json:"created_at"`
}
```

#### AccountAddress
```go
type AccountAddress struct {
    AccountID        string `json:"account_id"`
    TokenID          string `json:"token_id"`
    RecipientAddress string `json:"recipient_address"`
    RefName          string `json:"ref_name"`
    CreatedAt        string `json:"created_at"`
}
```

### Token Related

#### Token
```go
type Token struct {
    TokenID         string `json:"token_id"`
    Name            string `json:"name"`
    Symbol          string `json:"symbol"`
    ContractAddress string `json:"contract_address"`
    Decimals        int    `json:"decimals"`
    ChainID         int    `json:"chain_id"`
    ChainName       string `json:"chain_name"`
    ChainSymbol     string `json:"chain_symbol"`
    ExplorerURL     string `json:"explorer_url"`
    IconURL         string `json:"icon_url"`
    TokenType       string `json:"token_type"`
    IsActive        bool   `json:"is_active"`
    CurrencyType    string `json:"currency_type"`
    CreatedAt       string `json:"created_at"`
}
```

### Product Related

#### Product
```go
type Product struct {
    ProductID       string         `json:"product_id"`
    AccountID       string         `json:"account_id"`
    Name            string         `json:"name"`
    Description     string         `json:"description,omitempty"`
    Content         string         `json:"content"`
    Active          bool           `json:"active"`
    ProductTokens   []*ProductToken `json:"product_tokens"`
    CreatedAt       string         `json:"created_at"`
    TotalSaleCount  int64          `json:"total_sale_count"`
    TotalSaleAmount float64        `json:"total_sale_amount"`
}
```

#### ProductToken
```go
type ProductToken struct {
    ProductTokenID       string `json:"product_token_id"`
    ProductID            string `json:"product_id"`
    AccountID            string `json:"account_id"`
    TokenID              string `json:"token_id"`
    Price                string `json:"price"`
    RecipientAddress     string `json:"recipient_address"`
    PaymentRouterAddress string `json:"payment_router_address"`
    CreatedAt            string `json:"created_at"`
    ChainID              string `json:"chain_id"`
    ChainName            string `json:"chain_name"`
}
```

### Payment Related

#### Payment
```go
type Payment struct {
    PaymentID        string `json:"payment_id"`
    AccountID        string `json:"account_id"`
    TokenID          string `json:"token_id"`
    ProductID        string `json:"product_id"`
    ProductTokenID   string `json:"product_token_id"`
    Count            int    `json:"count"`
    Status           string `json:"status"`
    PayerEmail       string `json:"payer_email,omitempty"`
    CreatedAt        string `json:"created_at"`
    UpdatedAt        string `json:"updated_at"`
    PaidAt           string `json:"paid_at,omitempty"`
    ClosedAt         string `json:"closed_at,omitempty"`
    CloseReason      string `json:"close_reason,omitempty"`
    TransactionHash  string `json:"transaction_hash,omitempty"`
    BlockNumber      int64  `json:"block_number,omitempty"`
    GasUsed          int64  `json:"gas_used,omitempty"`
    GasPrice         string `json:"gas_price,omitempty"`
    TotalAmount      string `json:"total_amount"`
    FeeAmount        string `json:"fee_amount"`
    RecipientAmount  string `json:"recipient_amount"`
}
```

#### PaymentReceiver
```go
type PaymentReceiver struct {
    Type             string `json:"type"` // "fee" or "merchant"
    RecipientAddress string `json:"recipient_address"`
    Amount           string `json:"amount"` // wei format
    Rate             string `json:"rate"`   // percentage
}
```

#### ExternalCreatePaymentRequest
```go
type ExternalCreatePaymentRequest struct {
    ProductID      string `json:"product_id"`
    ProductTokenID string `json:"product_token_id"`
    Count          int    `json:"count"`
}
```

#### ExternalCreatePaymentResponse
```go
type ExternalCreatePaymentResponse struct {
    Message          string             `json:"message"`
    PaymentID        string             `json:"payment_id"`
    PayLink          string             `json:"pay_link"`
    ContractAddress  string             `json:"contract_address"`
    PaymentReceivers []*PaymentReceiver `json:"payment_receivers"`
    TokenAddress     string             `json:"token_address"`
    Decimals         int                `json:"decimals"`
}
```

## Error Handling

All methods in the SDK may return errors. Common error types include:

- **Network errors**: Request sending failed
- **Authentication errors**: JWT Token invalid or expired
- **Permission errors**: No access to specific resources
- **Parameter errors**: Invalid request parameters
- **Server errors**: Internal server errors

**Error handling example:**
```go
payment, err := client.GetPaymentByID("invalid-id")
if err != nil {
    if strings.Contains(err.Error(), "404") {
        fmt.Println("Payment not found")
    } else if strings.Contains(err.Error(), "401") {
        fmt.Println("Authentication failed, please check API Key")
    } else {
        fmt.Printf("Failed to get payment: %v\n", err)
    }
    return
}
```

## Authentication

The SDK uses JWT Token for authentication. The client automatically handles token acquisition and refresh:

1. Uses API Key to get access token during initialization
2. Automatically adds `Authorization: Bearer <token>` to request headers
3. Periodically refreshes tokens to maintain session validity

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/reddio-com/reddio-pay-sdk/go-sdk/client"
)

func main() {
    // Create client
    ctx := context.Background()
    client, err := client.NewSDKClient(ctx, "https://reddio-service-dev.reddio.com", "your-api-key")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Get account information
    account, err := client.GetAccountInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Welcome, %s!\n", account.Email)
    
    // Create product
    productReq := &client.CreateProductRequest{
        Name:             "Test Product",
        Description:      "This is a test product",
        Content:          "Product detailed content",
        TokenIDList:      []string{"token1"},
        Price:            "100.0",
        RecipientAddress: "0x1234567890abcdef",
    }
    
    productResp, err := client.CreateProduct(productReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Product created successfully, ID: %s\n", productResp.Product.ProductID)
    
    // Get product list
    products, err := client.ListProducts()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("You have %d products\n", len(products.Products))
    
    // Create external payment
    paymentReq := &client.ExternalCreatePaymentRequest{
        ProductID:      productResp.Product.ProductID,
        ProductTokenID: productResp.Product.ProductTokens[0].ProductTokenID,
        Count:          1,
    }
    
    paymentResp, err := client.ExternalCreatePayment(paymentReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("External payment created successfully, ID: %s\n", paymentResp.PaymentID)
    fmt.Printf("Payment link: %s\n", paymentResp.PayLink)
    
    // Get payment records
    payments, err := client.ListPaymentsByAccount(10, 0)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("You have %d payment records\n", payments.TotalCount)
}
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.
# Reddio Pay Go SDK Demonstration System

This is a Go-based demonstration system that showcases how to use various APIs from the `github.com/reddio-com/reddio-pay-sdk/go-sdk/client` SDK.

## Features

- **SDK API Demonstration**: Automatically demonstrates all SDK API usage methods on startup
- **Order Management**: Simple order creation and management functionality
- **Payment Integration**: Shows how to use the SDK to create payments and query status
- **RESTful API**: Provides standard REST API interfaces using Gorilla Mux

## Tech Stack

- Go 1.21
- Gorilla Mux Router
- database/sql + SQLite
- Reddio Pay Go SDK

## SDK Demonstration Features

The system automatically demonstrates the following SDK APIs on startup:

1. **List Products** (`ListProducts`)
2. **Create Product** (`CreateProduct`)
3. **Get Product Info** (`GetProduct`)
4. **Add Product Token** (`AddProductToken`)
5. **Create External Payment** (`ExternalCreatePayment`)
6. **Query Payment Status** (`GetPaymentByID`)

## Quick Start

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Run the Service

```bash
go run main.go
```

Or compile and run:

```bash
go build -o order-system main.go
./order-system
```

The service will start on `http://localhost:8080`

### 3. Environment Variables (Optional)

```bash
export REDDIO_API_KEY="your_api_key"
export REDDIO_URL="https://reddio-service-prod.reddio.com"
```

## API Endpoints

### Basic Information

- **Base URL**: `http://localhost:8080/api`
- **Content-Type**: `application/json`

### 1. Create Order

```http
POST /api/orders
Content-Type: application/json

{
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "product_id": "prod_123",
  "product_token_id": "token_456",
  "quantity": 1
}
```

**Response Example:**
```json
{
  "id": 1,
  "order_number": "ORD00000001",
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "product_id": "prod_123",
  "product_token_id": "token_456",
  "quantity": 1,
  "total_amount": "100.00",
  "status": "pending",
  "reddio_payment_id": "pay_789",
  "reddio_pay_link": "https://pay.reddio.com/pay/pay_789",
  "reddio_status": "created",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 2. Get Order Details

```http
GET /api/orders/{id}
```

**Response Example:**
```json
{
  "id": 1,
  "order_number": "ORD00000001",
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "product_id": "prod_123",
  "product_token_id": "token_456",
  "quantity": 1,
  "total_amount": "100.00",
  "status": "pending",
  "reddio_payment_id": "pay_789",
  "reddio_pay_link": "https://pay.reddio.com/pay/pay_789",
  "reddio_status": "created",
  "transaction_hash": "",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "paid_at": null
}
```

### 3. Get Order Status

```http
GET /api/orders/{id}/status
```

**Response Example:**
```json
{
  "order_number": "ORD00000001",
  "status": "pending",
  "reddio_status": "created",
  "reddio_pay_link": "https://pay.reddio.com/pay/pay_789",
  "transaction_hash": "",
  "paid_at": null
}
```

### 4. Check Payment Status

```http
POST /api/orders/{id}/check-payment
```

**Response Example:**
```json
{
  "message": "Payment status updated",
  "order_number": "ORD00000001",
  "status": "paid",
  "reddio_status": "paid",
  "transaction_hash": "0x1234567890abcdef",
  "paid_at": "2024-01-01T01:00:00Z"
}
```

### 5. Get Order List

```http
GET /api/orders?page=1&limit=10&status=pending
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `status`: Order status filter (optional)

**Response Example:**
```json
{
  "orders": [
    {
      "id": 1,
      "order_number": "ORD00000001",
      "customer_name": "John Doe",
      "customer_email": "john@example.com",
      "product_id": "prod_123",
      "product_token_id": "token_456",
      "quantity": 1,
      "total_amount": "100.00",
      "status": "pending",
      "reddio_payment_id": "pay_789",
      "reddio_pay_link": "https://pay.reddio.com/pay/pay_789",
      "reddio_status": "created",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

## Order Status

- `pending`: Awaiting payment
- `paid`: Payment completed
- `cancelled`: Order cancelled
- `expired`: Order expired
- `failed`: Creation failed

## Usage Flow

1. **Create Order**: Call `/api/orders` to create an order and get a Reddio Pay payment link
2. **User Payment**: User visits the returned `reddio_pay_link` to complete payment
3. **Check Status**: Call `/api/orders/{id}/check-payment` to update payment status
4. **View Order**: Call `/api/orders/{id}` to view order details

## Database

The system uses the native `database/sql` package to operate on SQLite database, with the database file being `orders.db`.

### Database Table Structure

- **orders**: Order table (for demonstrating order management functionality)

## Logging

The system outputs detailed log information including:
- SDK call logs
- Order creation logs
- Payment status update logs
- Error messages

## Error Handling

All APIs return standard error responses:

```json
{
  "error": "Error description"
}
```

Common HTTP status codes:
- `200`: Success
- `201`: Created successfully
- `400`: Invalid request parameters
- `404`: Resource not found
- `500`: Internal server error
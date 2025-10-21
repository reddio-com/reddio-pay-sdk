# Reddio Pay TypeScript SDK

TypeScript SDK for Reddio Pay services - A complete payment solution for crypto transactions.

## Features

- 🔐 Automatic JWT authentication and token refresh
- 🛍️ Product management (create, list, manage digital products)
- 💳 Payment processing (create orders, query status)
- 🪙 Token support (multi-chain crypto tokens)
- 📊 Payment analytics and reporting
- 🔔 Webhook notifications
- 📘 Full TypeScript support with type definitions

## Installation

```bash
npm install @reddio/pay-sdk
# or
yarn add @reddio/pay-sdk
```

## Quick Start

```typescript
import { ReddioClient } from '@reddio/pay-sdk';

const client = new ReddioClient({
  baseURL: 'https://reddio-service-prod.reddio.com',
  apiKey: 'your-api-key'
});

// Initialize the client
await client.initialize();

// Get product list
const products = await client.product.listProducts();
console.log(`Found ${products.length} products`);
```

## API Reference

### Product Management
- `listProducts()` - Get all products
- `createProduct()` - Create a new product
- `getProduct(id)` - Get product by ID
- `addProductToken()` - Add token support to product

### Payment Processing
- `createPayment()` - Create new payment order
- `getPayment(id)` - Get payment status
- `listPayments()` - List payment history

## License

MIT

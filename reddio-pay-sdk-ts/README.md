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
- 🌐 Support for production and development environments

## Installation

```bash
npm install reddio-pay-sdk
# or
yarn add reddio-pay-sdk
```

## Quick Start

```typescript
import { ReddioClient } from 'reddio-pay-sdk';

// Create client for production environment
const client = new ReddioClient({
  apiKey: 'your-api-key',
  environment: 'prod' // or 'dev' for development
});

// Initialize the client (authenticate)
await client.initialize();

// Get product list
const products = await client.product.listProducts();
console.log(`Found ${products.length} products`);

// Clean up when done
client.destroy();
```

## API Reference

### Client Initialization

```typescript
// Standard initialization
const client = new ReddioClient({
  apiKey: 'your-api-key',
  environment: 'prod', // 'prod' or 'dev'
  timeout: 30000 // optional, default 30 seconds
});

// Factory methods
const prodClient = ReddioClient.createProd('your-api-key');
const devClient = ReddioClient.createDev('your-api-key');
const customClient = ReddioClient.create('https://custom-url.com', 'your-api-key');
```

### Product Management

```typescript
// List all products
const products = await client.product.listProducts();

// Get product by ID
const product = await client.product.getProduct('product_id');

// Create new product
const newProduct = await client.product.createProduct({
  name: 'Premium Plan',
  description: 'Full access to all features',
  amount: '99.00',
  currency: 'USD'
});

// Add token support to product
await client.product.addProductToken('product_id', 'token_id');
```

### Payment Processing

```typescript
// Create payment order
const payment = await client.payment.createPayment({
  productId: 'product_id',
  receiverEmail: 'customer@example.com',
  amount: '99.00',
  currency: 'USD',
  description: 'Order #12345',
  callbackUrl: 'https://your-site.com/callback'
});

// Get payment status
const status = await client.payment.getPayment('payment_id');

// List payment history
const payments = await client.payment.listPayments({
  limit: 10,
  offset: 0
});

// Send payment success notification
await client.payment.sendPaymentSuccessNotification({
  paymentId: 'payment_id',
  customMessage: 'Thank you for your purchase!'
});
```

### Token Management

```typescript
// List available tokens
const tokens = await client.token.listTokens();

// Get token details
const token = await client.token.getToken('token_id');

// Create new token
const newToken = await client.token.createToken({
  name: 'USDT',
  symbol: 'USDT',
  chainId: 1,
  contractAddress: '0x...'
});
```

### Webhook Configuration

```typescript
// Update webhook URL
await client.account.updateWebhook('https://your-site.com/webhook');
```

## Error Handling

```typescript
import { ReddioPayError, AuthenticationError, NetworkError } from 'reddio-pay-sdk';

try {
  await client.payment.createPayment(data);
} catch (error) {
  if (error instanceof AuthenticationError) {
    console.error('Authentication failed:', error.message);
  } else if (error instanceof NetworkError) {
    console.error('Network error:', error.message);
  } else {
    console.error('Unknown error:', error);
  }
}
```

## Environment Variables

Enable debug logging:
```bash
export LOG_LEVEL=debug
```

## Examples

See the [examples](./examples) directory for more detailed usage examples.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- GitHub Issues: [https://github.com/reddio-com/reddio-pay-sdk/issues](https://github.com/reddio-com/reddio-pay-sdk/issues)
- Documentation: [https://docs.reddio.com](https://docs.reddio.com)
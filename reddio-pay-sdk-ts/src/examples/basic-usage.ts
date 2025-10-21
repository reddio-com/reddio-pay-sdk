import { ReddioClient } from '../index';

/**
 * Basic usage example demonstrating how to list products
 */
async function basicUsageExample() {
  console.log('=== Reddio Pay SDK Basic Usage Example ===\n');

  // Create client instance
  const client = new ReddioClient({
    baseURL: 'https://reddio-service-dev.reddio.com', // Use dev environment for testing
    apiKey: 'your-api-key-here' // Replace with your actual API key
  });

  try {
    // Initialize the client (authenticate)
    console.log('1. Initializing Reddio client...');
    await client.initialize();
    console.log('✅ Client initialized successfully\n');

    // Test 1: List supported tokens
    console.log('2. Fetching supported tokens...');
    const tokens = await client.token.listTokens();
    console.log(`✅ Found ${tokens.length} supported tokens:`);
    tokens.slice(0, 3).forEach((token, index) => {
      console.log(`   ${index + 1}. ${token.name} (${token.symbol}) - Chain: ${token.chainName}`);
    });
    if (tokens.length > 3) {
      console.log(`   ... and ${tokens.length - 3} more tokens`);
    }
    console.log();

    // Test 2: List products (main functionality we want to test)
    console.log('3. Fetching product list...');
    const products = await client.product.listProducts();
    console.log(`✅ Found ${products.length} products:`);
    
    if (products.length > 0) {
      products.forEach((product, index) => {
        console.log(`   ${index + 1}. Product: ${product.name}`);
        console.log(`      ID: ${product.productId}`);
        console.log(`      Active: ${product.active}`);
        console.log(`      Sales Count: ${product.totalSaleCount}`);
        console.log(`      Sales Amount: ${product.totalSaleAmount}`);
        console.log(`      Supported Tokens: ${product.productTokens.length}`);
        console.log();
      });
    } else {
      console.log('   No products found. You can create one using client.product.createProduct()');
      console.log();
    }

    console.log('🎉 All tests completed successfully!');

  } catch (error) {
    console.error('❌ Error occurred:', error);
    
    if (error instanceof Error) {
      console.error('Error message:', error.message);
    }
    
    // Common error scenarios and solutions
    console.log('\n💡 Common solutions:');
    console.log('- Ensure your API key is correct and active');
    console.log('- Check if the base URL is correct (dev/prod environment)');
    console.log('- Verify your network connection');
    console.log('- Make sure your account has the necessary permissions');
    
  } finally {
    // Cleanup
    console.log('\n4. Cleaning up...');
    client.destroy();
    console.log('✅ Client cleaned up');
  }
}

/**
 * Advanced example showing product creation
 */
async function advancedExample() {
  console.log('\n=== Advanced Example: Create Product ===\n');

  const client = new ReddioClient({
    baseURL: 'https://reddio-service-dev.reddio.com',
    apiKey: 'your-api-key-here' // Replace with your actual API key
  });

  try {
    await client.initialize();
    
    // Get supported tokens first
    const tokens = await client.token.listTokens();
    if (tokens.length === 0) {
      throw new Error('No supported tokens found');
    }

    const firstToken = tokens[0];
    console.log(`Using token: ${firstToken.name} (${firstToken.symbol})`);

    // Create a demo product
    console.log('\n1. Creating demo product...');
    const newProduct = await client.product.createProduct({
      name: 'Demo Digital Product',
      description: 'A sample product created via TypeScript SDK',
      content: 'This is the content of the digital product',
      tokenIdList: [firstToken.tokenId],
      price: '1000000000000000000', // 1 token in wei format
      recipientAddress: '0x1234567890abcdef1234567890abcdef12345678' // Replace with actual address
    });

    console.log(`✅ Product created: ${newProduct.name} (ID: ${newProduct.productId})`);

    // List products again to see the new one
    console.log('\n2. Listing products after creation...');
    const updatedProducts = await client.product.listProducts();
    console.log(`✅ Now have ${updatedProducts.length} products total`);

  } catch (error) {
    console.error('❌ Advanced example failed:', error);
  } finally {
    client.destroy();
  }
}

// Run the examples
async function main() {
  // Check if API key is provided
  if (process.argv.includes('--api-key')) {
    const apiKeyIndex = process.argv.indexOf('--api-key') + 1;
    if (apiKeyIndex < process.argv.length) {
      // Replace the placeholder with actual API key
      console.log('Using provided API key for testing');
    }
  } else {
    console.log('⚠️  To run with real API calls, provide --api-key YOUR_API_KEY');
    console.log('Currently using placeholder API key (will fail authentication)');
  }

  // Run basic example
  await basicUsageExample();
  
  // Uncomment to run advanced example
  // await advancedExample();
}

// Run if this file is executed directly
if (require.main === module) {
  main().catch(console.error);
}

export { basicUsageExample, advancedExample };

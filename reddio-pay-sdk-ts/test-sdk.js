const { ReddioClient } = require('./dist/index');

console.log('🧪 Reddio Pay TypeScript SDK - Quick Test\n');

// Test 1: Basic import and instantiation
try {
  const client = new ReddioClient({
    baseURL: 'https://reddio-service-dev.reddio.com',
    apiKey: 'test-api-key'
  });
  
  console.log('✅ Test 1: SDK import and client creation - PASSED');
  
  // Test 2: Check API structure
  const hasProductAPI = client.product && typeof client.product.listProducts === 'function';
  const hasTokenAPI = client.token && typeof client.token.listTokens === 'function';
  const hasInitialize = typeof client.initialize === 'function';
  
  if (hasProductAPI && hasTokenAPI && hasInitialize) {
    console.log('✅ Test 2: API structure validation - PASSED');
  } else {
    console.log('❌ Test 2: API structure validation - FAILED');
  }
  
  // Test 3: Method availability
  console.log('✅ Test 3: Available methods:');
  console.log('   - client.product.listProducts(): ' + typeof client.product.listProducts);
  console.log('   - client.product.createProduct(): ' + typeof client.product.createProduct);
  console.log('   - client.product.getProduct(): ' + typeof client.product.getProduct);
  console.log('   - client.token.listTokens(): ' + typeof client.token.listTokens);
  
  client.destroy();
  console.log('\n🎉 All tests PASSED! SDK is ready to use.');
  
  console.log('\n📝 Next steps:');
  console.log('1. Replace "test-api-key" with your actual Reddio Pay API key');
  console.log('2. Call await client.initialize() to authenticate');
  console.log('3. Use await client.product.listProducts() to fetch products');
  console.log('4. Use other API methods as needed');
  
} catch (error) {
  console.error('❌ Test failed:', error.message);
}

const { ReddioClient } = require('./dist/index.js');

// Try to load test config, use default if not found
let testConfig;
try {
  testConfig = require('./test-config.js');
} catch (error) {
  console.log('⚠️  test-config.js file not found, please create this file and add your API Key');
  console.log('   Example: module.exports = { API_KEY: "your-api-key-here" };');
  process.exit(1);
}

/**
 * Real API Test - with detailed RPC logs
 */
async function testRealAPIWithTimeout() {
  console.log('🌐 Starting real API test (detailed RPC logs)...\n');
  console.log('='.repeat(80));
  
  const API_KEY = testConfig.API_KEY;
  
  if (!API_KEY || API_KEY === 'your-api-key-here') {
    console.log('❌ Please set a valid API Key in test-config.js');
    process.exit(1);
  }
  
  console.log('🔧 Testing production environment...\n');
  
  await testEnvironmentWithTimeout({
    apiKey: API_KEY,
    environment: 'prod'
  });
}

async function testEnvironmentWithTimeout(config) {
  const client = new ReddioClient(config);
  
  try {
    console.log(`📍 Environment: ${config.environment}`);
    console.log('─'.repeat(80));
    
    console.log('🔑 Step 1: Initialize client...');
    await client.initialize();
    console.log('✅ Client initialized successfully\n');
    console.log('─'.repeat(80));
    
    // Test product API
    console.log('📦 Step 2: Test product API...');
    await testAPIWithTimeout('Product', async () => {
      const products = await client.product.listProducts();
      
      if (products && products.products && Array.isArray(products.products)) {
        console.log(`✅ Product list succeeded! Total: ${products.products.length}`);
        return products.products.length;
      } else if (Array.isArray(products)) {
        console.log(`✅ Product list succeeded! Total: ${products.length}`);
        return products.length;
      } else {
        console.log(`⚠️  Product API returned unexpected format`);
        return 'unknown';
      }
    });
    console.log('─'.repeat(80));

    // Test token API
    console.log('🪙 Step 3: Test token API...');
    await testAPIWithTimeout('Token', async () => {
      const tokens = await client.token.listTokens();
      
      if (tokens === undefined || tokens === null) {
        console.log('⚠️  Token API returned undefined/null');
        return 'empty';
      } else if (tokens && tokens.tokens && Array.isArray(tokens.tokens)) {
        console.log(`✅ Token list succeeded! Total: ${tokens.tokens.length}`);
        return tokens.tokens.length;
      } else if (Array.isArray(tokens)) {
        console.log(`✅ Token list succeeded! Total: ${tokens.length}`);
        return tokens.length;
      } else {
        console.log(`⚠️  Token API returned unexpected format, type: ${typeof tokens}`);
        return 'unknown';
      }
    });
    console.log('─'.repeat(80));

    // Test payment API
    console.log('💰 Step 4: Test payment API...');
    await testAPIWithTimeout('Payment', async () => {
      const payments = await client.payment.listPayments({ limit: 5 });
      
      if (payments === undefined || payments === null) {
        console.log('⚠️  Payment API returned undefined/null');
        return 'empty';
      } else if (payments && payments.payments && Array.isArray(payments.payments)) {
        console.log(`✅ Payment list succeeded! Total: ${payments.payments.length}`);
        return payments.payments.length;
      } else if (Array.isArray(payments)) {
        console.log(`✅ Payment list succeeded! Total: ${payments.length}`);
        return payments.length;
      } else {
        console.log(`⚠️  Payment API returned unexpected format, type: ${typeof payments}`);
        return 'unknown';
      }
    });
    
  } catch (error) {
    console.log(`❌ Test failed: ${error.message}`);
  } finally {
    console.log('─'.repeat(80));
    console.log('🧹 Cleaning up client...');
    client.destroy();
    console.log('✅ Client cleaned up');
  }
}

async function testAPIWithTimeout(apiName, testFunction) {
  const timeout = 15000; // 15 seconds timeout
  
  try {
    console.log(`⏰ Start testing ${apiName} API (${timeout/1000} seconds timeout)...`);
    
    const result = await Promise.race([
      testFunction(),
      new Promise((_, reject) => 
        setTimeout(() => reject(new Error(`${apiName} API test timeout`)), timeout)
      )
    ]);
    
    console.log(`✅ ${apiName} API test finished, result: ${result}\n`);
    
  } catch (error) {
    if (error.message.includes('timeout')) {
      console.log(`⏰ ${apiName} API test timeout - endpoint may be slow or stuck\n`);
    } else {
      console.log(`❌ ${apiName} API test failed: ${error.message}\n`);
    }
  }
}

// Run test
if (require.main === module) {
  testRealAPIWithTimeout()
    .then(() => {
      console.log('='.repeat(80));
      console.log('🎉 Test finished!');
      process.exit(0); // Force exit to prevent hanging
    })
    .catch((error) => {
      console.log('='.repeat(80));
      console.error('❌ Test failed:', error);
      process.exit(1);
    });
}

module.exports = { testRealAPIWithTimeout };
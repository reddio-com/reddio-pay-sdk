const { ReddioClient, PaymentApi } = require('./dist/index.js');

console.log('🧪 Payment API Mock Test...\n');

// Mock HTTP Client
class MockHttpClient {
  constructor(config) {
    this.config = config;
    console.log('📡 Mock HTTP Client initialized:', config.baseURL);
  }

  async initialize() {
    console.log('🔐 Mock client authentication succeeded');
    return Promise.resolve();
  }

  async post(url, data) {
    console.log('📤 Mock POST:', url);
    console.log('   Data:', JSON.stringify(data, null, 2));
    
    if (url === '/external/payments') {
      return {
        data: {
          payment: {
            id: 'payment_mock_123456',
            accountId: 'account_mock_001',
            productId: data.productId,
            receiverId: 'receiver_mock_001',
            amount: data.amount,
            currency: data.currency,
            status: 'pending',
            description: data.description,
            metadata: data.metadata,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
          },
          paymentUrl: 'https://pay.reddio.com/payment/payment_mock_123456'
        }
      };
    }
    
    if (url === '/external/payments/success/notify') {
      return {
        data: {
          success: true,
          message: 'Payment success notification sent successfully'
        }
      };
    }
    
    return { data: {} };
  }

  async get(url) {
    console.log('📥 Mock GET:', url);
    
    if (url.includes('/payments/payment_mock_123456')) {
      return {
        data: {
          id: 'payment_mock_123456',
          accountId: 'account_mock_001',
          productId: 'test_product_123',
          receiverId: 'receiver_mock_001',
          amount: '10.00',
          currency: 'USD',
          status: 'completed',
          txHash: '0x123456789abcdef',
          description: 'Payment API test order',
          metadata: { testCase: 'payment_api_test' },
          createdAt: '2024-10-22T10:00:00Z',
          updatedAt: '2024-10-22T10:05:00Z',
          receiver: {
            id: 'receiver_mock_001',
            email: 'test@example.com',
            firstName: 'Test',
            lastName: 'User',
            createdAt: '2024-10-20T10:00:00Z',
            updatedAt: '2024-10-20T10:00:00Z'
          }
        }
      };
    }
    
    if (url.includes('/payments/product/')) {
      return {
        data: {
          payments: [
            {
              id: 'payment_mock_123456',
              amount: '10.00',
              currency: 'USD',
              status: 'completed',
              createdAt: '2024-10-22T10:00:00Z'
            }
          ]
        }
      };
    }
    
    if (url.includes('/payments/list')) {
      return {
        data: {
          payments: [
            {
              id: 'payment_mock_123456',
              amount: '10.00',
              currency: 'USD',
              status: 'completed',
              createdAt: '2024-10-22T10:00:00Z'
            }
          ],
          pagination: {
            totalCount: 1,
            totalPages: 1,
            currentPage: 1,
            pageSize: 5
          }
        }
      };
    }
    
    return { data: {} };
  }

  destroy() {
    console.log('🧹 Mock client resources cleaned up');
  }
}

async function runMockTest() {
  try {
    console.log('1️⃣ Creating Mock Payment API client...');
    const mockHttpClient = new MockHttpClient({
      baseURL: 'https://api.reddio.com',
      apiKey: 'mock-api-key'
    });
    
    const paymentApi = new PaymentApi(mockHttpClient);
    await mockHttpClient.initialize();
    console.log('✅ Mock client created successfully\n');

    // Test data
    const testData = {
      productId: 'test_product_123',
      receiverEmail: 'test@example.com',
      amount: '10.00',
      currency: 'USD',
      description: 'Payment API test order',
      metadata: {
        testCase: 'payment_api_test',
        timestamp: new Date().toISOString()
      },
      callbackUrl: 'https://example.com/callback'
    };

    console.log('2️⃣ Testing create payment order...');
    const createdPayment = await paymentApi.createPayment(testData);
    console.log('✅ Payment order created successfully:');
    console.log('   Payment ID:', createdPayment.payment.id);
    console.log('   Payment URL:', createdPayment.payLink);
    console.log('   Status:', createdPayment.payment.status);

    console.log('\n3️⃣ Testing get payment details...');
    const paymentDetails = await paymentApi.getPayment(createdPayment.payment.id);
    console.log('✅ Get payment details succeeded:');
    console.log('   Payment ID:', paymentDetails.id);
    console.log('   Status:', paymentDetails.status);
    console.log('   Receiver Email:', paymentDetails.receiver?.email);

    console.log('\n4️⃣ Testing get payments by product ID...');
    const productPayments = await paymentApi.listPaymentsByProduct(testData.productId);
    console.log('✅ Get payments by product ID succeeded:');
    console.log('   Number of payments:', productPayments.payments.length);

    console.log('\n5️⃣ Testing list all payments with pagination...');
    const allPayments = await paymentApi.listPayments({ limit: 5, offset: 0 });
    console.log('✅ Get all payments succeeded:');
    console.log('   Number of payments:', allPayments.payments.length);
    console.log('   Total count:', allPayments.pagination.totalCount);

    console.log('\n6️⃣ Testing send payment success notification...');
    const notifyResult = await paymentApi.sendPaymentSuccessNotification({
      paymentId: createdPayment.payment.id,
      customMessage: 'Mock test notification message'
    });
    console.log('✅ Notification sent successfully:');
    console.log('   Success:', notifyResult.success);
    console.log('   Message:', notifyResult.message);

    mockHttpClient.destroy();

    console.log('\n🎉 All Payment API method tests passed!');
    console.log('\n📋 Test results:');
    console.log('   ✅ createPayment - Create payment order');
    console.log('   ✅ getPayment - Get payment details');
    console.log('   ✅ listPaymentsByProduct - Get payments by product ID');
    console.log('   ✅ listPayments - List all payments with pagination');
    console.log('   ✅ sendPaymentSuccessNotification - Send payment success notification');

  } catch (error) {
    console.error('❌ Mock test failed:', error.message);
    console.error(error.stack);
  }
}

runMockTest();
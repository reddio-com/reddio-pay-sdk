import { ReddioClient } from '../client/reddio-client';

async function paymentExample() {
  const client = new ReddioClient({
    baseURL: 'https://api.reddio.com',
    apiKey: 'your-api-key',
    timeout: 30000,
  });

  try {
    // 初始化客户端
    await client.initialize();

    // 1. 创建支付订单
  const paymentRequest = {
    product_id: 'prod_example_123',
    product_token_id: 'token_example_456',
    count: 1
  };

    const createdPayment = await client.payment.createPayment(paymentRequest);
    console.log('Payment created:', createdPayment);
    console.log('Payment URL:', createdPayment.pay_link);

    // 2. 获取支付详情
    const paymentDetails = await client.payment.getPayment(createdPayment.payment_id);
    console.log('Payment details:', paymentDetails);

    // 3. 根据产品ID获取支付列表
    const productPayments = await client.payment.listPaymentsByProduct('prod_example_123');
    console.log('Product payments:', productPayments);

    // 4. 分页获取所有支付记录
    const allPayments = await client.payment.listPayments({
      limit: 10,
      offset: 0
    });
    console.log('All payments:', allPayments);

    // 5. 发送支付成功通知
    const notifyRequest = {
      paymentId: createdPayment.payment_id,
      customMessage: '感谢您的购买！'
    };

    const notifyResult = await client.payment.sendPaymentSuccessNotification(notifyRequest);
    console.log('Notification sent:', notifyResult);

  } catch (error) {
    console.error('Payment API error:', error);
  } finally {
    // 清理资源
    client.destroy();
  }
}

// 导出示例函数
export { paymentExample };

// 如果直接运行此文件，执行示例
if (require.main === module) {
  paymentExample().catch(console.error);
}

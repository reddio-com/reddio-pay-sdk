const { ReddioClient } = require('./dist/index.js');
const testConfig = require('./test-config.js');

async function testUpdateWebhook() {
  console.log('--- Test Webhook Start ---');
  console.log('API_KEY:', testConfig.API_KEY);
  console.log('Environment:', 'prod');
  const webhookUrl = 'https://example.com/webhook';
  console.log('Webhook URL:', webhookUrl);

  const client = new ReddioClient({
    apiKey: testConfig.API_KEY,
    environment: 'prod'
  });

  try {
    console.log('Initializing client...');
    await client.initialize();
    console.log('Client initialized.');

    // 检查 HttpClient 是否已保存 token
    if (client.httpClient && client.httpClient.accessToken) {
      console.log('Access Token:', client.httpClient.accessToken);
    } else {
      console.log('No access token found after initialization.');
    }

    console.log('Calling updateWebhook...');
    const resp = await client.account.updateWebhook(webhookUrl);
    console.log('Webhook update response:', resp);

    // 检查响应内容
    if (resp && resp.message) {
      console.log('Webhook update message:', resp.message);
    } else {
      console.log('No message in webhook update response.');
    }
  } catch (error) {
    console.error('Webhook update failed:', error);
    if (error.response) {
      console.error('Error response data:', error.response.data);
      console.error('Error status:', error.response.status);
      console.error('Error headers:', error.response.headers);
    }
    if (error.statusCode) {
      console.error('Error statusCode:', error.statusCode);
    }
    if (error.code) {
      console.error('Error code:', error.code);
    }
  }
  console.log('--- Test Webhook End ---');
}

testUpdateWebhook();
import { HttpClient } from '../utils/http-client';
import { UpdateWebhookRequest, UpdateWebhookResponse } from '../types/common';

export class AccountApi {
  constructor(private httpClient: HttpClient) {}

  async updateWebhook(webhookUrl: string): Promise<UpdateWebhookResponse> {
    const reqBody: UpdateWebhookRequest = { webhook: webhookUrl };
    return await this.httpClient.put<UpdateWebhookResponse>('/accounts/webhook', reqBody);
  }
}
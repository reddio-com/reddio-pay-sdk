import { HttpClient } from '../utils/http-client';
import { PaginationOptions } from '../types/common';
import {
  Payment,
  ExternalCreatePaymentRequest,
  ExternalCreatePaymentResponse,
  ListPaymentsResponse,
  ListPaymentsResponseWithPagination,
  ExternalSendNotifyForPaymentSuccessRequest,
  SendNotifyForPaymentSuccessResponse
} from '../types/payment';


export class PaymentApi {
  constructor(private httpClient: HttpClient) {}

  async createPayment(request: ExternalCreatePaymentRequest): Promise<ExternalCreatePaymentResponse> {
    const response = await this.httpClient.post('/external/payments', request);
    return response;
  }


  async getPayment(paymentId: string): Promise<Payment> {
    const response = await this.httpClient.get(`/payments/${paymentId}`);
    return response;
  }


  async listPaymentsByProduct(productId: string): Promise<ListPaymentsResponse> {
    const response = await this.httpClient.get(`/payments/product/${productId}`);
    return response;
  }

  async listPayments(options: PaginationOptions = {}): Promise<ListPaymentsResponseWithPagination> {
    const params = new URLSearchParams();
    if (options.limit) {
      params.append('limit', options.limit.toString());
    }
    if (options.offset !== undefined) {
      params.append('offset', options.offset.toString());
    }
    
    const url = `/payments/list${params.toString() ? '?' + params.toString() : ''}`;
    const response = await this.httpClient.get(url);
    return response;
  }

 
  async sendPaymentSuccessNotification(request: ExternalSendNotifyForPaymentSuccessRequest): Promise<SendNotifyForPaymentSuccessResponse> {
    const response = await this.httpClient.post('/external/payments/success/notify', request);
    return response;
  }
}

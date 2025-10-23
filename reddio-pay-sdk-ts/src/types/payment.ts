import { PaginationOptions } from './common';


export interface PaymentReceiver {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  createdAt: string;
  updatedAt: string;
}


export interface PaginationMeta {
  totalCount: number;
  totalPages: number;
  currentPage: number;
  pageSize: number;
}


export interface Payment {
  id: string;
  accountId: string;
  productId: string;
  receiverId: string;
  amount: string;
  currency: string;
  status: string;
  txHash?: string;
  description?: string;
  metadata?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
  receiver?: PaymentReceiver;
}


export interface PaymentReceiver {
  type: string;               // "fee" or "merchant"
  recipient_address: string;
  amount: string;             // wei format
  rate: string;               // percentage
}


export interface ExternalCreatePaymentRequest {
  product_id: string;
  product_token_id: string;
  count: number;
}


export interface ExternalCreatePaymentResponse {
  message: string;
  payment_id: string;
  pay_link: string;
  contract_address: string;
  payment_receivers: PaymentReceiver[];
  token_address: string;
  decimals: number;
}


export interface ListPaymentsResponse {
  payments: Payment[];
}


export interface ListPaymentsResponseWithPagination {
  payments: Payment[];
  pagination: PaginationMeta;
}


export interface ExternalSendNotifyForPaymentSuccessRequest {
  paymentId: string;
  customMessage?: string;
}


export interface SendNotifyForPaymentSuccessResponse {
  success: boolean;
  message: string;
}

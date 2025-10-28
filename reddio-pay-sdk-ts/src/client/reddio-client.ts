import { HttpClient } from '../utils/http-client';
import { ProductApi } from './product-api';
import { TokenApi } from './token-api';
import { PaymentApi } from './payment-api';
import { AccountApi } from './account-api';

import { ClientConfig, resolveClientConfig } from '../types/common';

/**
 * Main Reddio Pay SDK client
 */
export class ReddioClient {
  private httpClient: HttpClient;
  public readonly product: ProductApi;
  public readonly token: TokenApi;
  public readonly payment: PaymentApi;
  public readonly account: AccountApi; 

  constructor(config: ClientConfig) {
    const resolvedConfig = resolveClientConfig(config);
    
    this.httpClient = new HttpClient(resolvedConfig);
    this.product = new ProductApi(this.httpClient);
    this.token = new TokenApi(this.httpClient);
    this.payment = new PaymentApi(this.httpClient);
    this.account = new AccountApi(this.httpClient); 
  }

  /**
   * Static factory method: create client for production environment
   */
  static createProd(apiKey: string, options?: Partial<ClientConfig>): ReddioClient {
    return new ReddioClient({
      apiKey,
      environment: 'prod',
      ...options
    });
  }

  /**
   * Static factory method: create client for development environment
   */
  static createDev(apiKey: string, options?: Partial<ClientConfig>): ReddioClient {
    return new ReddioClient({
      apiKey,
      environment: 'dev',
      ...options
    });
  }

  /**
   * Static factory method: create client for custom environment
   */
  static create(baseURL: string, apiKey: string, options?: Partial<ClientConfig>): ReddioClient {
    return new ReddioClient({
      baseURL,
      apiKey,
      ...options
    });
  }

  /**
   * Initialize the client (authenticate and setup)
   */
  async initialize(): Promise<void> {
    await this.httpClient.initialize();
  }

  /**
   * Cleanup and destroy the client
   */
  destroy(): void {
    this.httpClient.destroy();
  }

}
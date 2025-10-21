import { HttpClient } from '../utils/http-client';
import { ProductApi } from './product-api';
import { TokenApi } from './token-api';
import { ClientConfig } from '../types';

/**
 * Main Reddio Pay SDK client
 */
export class ReddioClient {
  private httpClient: HttpClient;
  public readonly product: ProductApi;
  public readonly token: TokenApi;

  constructor(config: ClientConfig) {
    this.httpClient = new HttpClient(config);
    this.product = new ProductApi(this.httpClient);
    this.token = new TokenApi(this.httpClient);
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

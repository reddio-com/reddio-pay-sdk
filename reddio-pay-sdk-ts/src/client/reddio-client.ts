import { HttpClient } from '../utils/http-client';
import { ProductApi } from './product-api';
import { TokenApi } from './token-api';
import { PaymentApi } from './payment-api';
import { ClientConfig, resolveClientConfig } from '../types/common';

/**
 * Main Reddio Pay SDK client
 */
export class ReddioClient {
  private httpClient: HttpClient;
  public readonly product: ProductApi;
  public readonly token: TokenApi;
  public readonly payment: PaymentApi;

  constructor(config: ClientConfig) {
    // 解析配置，应用默认值
    const resolvedConfig = resolveClientConfig(config);
    
    this.httpClient = new HttpClient(resolvedConfig);
    this.product = new ProductApi(this.httpClient);
    this.token = new TokenApi(this.httpClient);
    this.payment = new PaymentApi(this.httpClient);
  }

  /**
   * 静态工厂方法：创建生产环境客户端
   */
  static createProd(apiKey: string, options?: Partial<ClientConfig>): ReddioClient {
    return new ReddioClient({
      apiKey,
      environment: 'prod',
      ...options
    });
  }

  /**
   * 静态工厂方法：创建开发环境客户端
   */
  static createDev(apiKey: string, options?: Partial<ClientConfig>): ReddioClient {
    return new ReddioClient({
      apiKey,
      environment: 'dev',
      ...options
    });
  }

  /**
   * 静态工厂方法：创建自定义环境客户端
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
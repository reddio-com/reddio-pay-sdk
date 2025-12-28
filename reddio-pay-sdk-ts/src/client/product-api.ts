import { HttpClient } from '../utils/http-client';
import {
  Product,
  ListProductsResponse,
  CreateProductRequest,
  CreateProductResponse,
  AddProductTokenRequest,
  AddProductTokenResponse,
  GetProductTokenStatusResponse
} from '../types';

/**
 * Product API client for managing digital products
 */
export class ProductApi {
  constructor(private httpClient: HttpClient) {}

  /**
   * Get all products for the authenticated account
   */
  async listProducts(): Promise<Product[]> {
    const response = await this.httpClient.get<ListProductsResponse>('/products');
    return response.products;
  }

  /**
   * Create a new product
   */
  async createProduct(request: CreateProductRequest): Promise<Product> {
    const response = await this.httpClient.post<CreateProductResponse>('/products', request);
    return response.product;
  }

  /**
   * Get product details by ID
   */
  async getProduct(productId: string): Promise<Product> {
    return await this.httpClient.get<Product>(`/products/${productId}`);
  }

  /**
   * Add a token to an existing product
   */
  async addProductToken(productId: string, request: AddProductTokenRequest): Promise<AddProductTokenResponse> {
    return await this.httpClient.post<AddProductTokenResponse>(
      `/products/${productId}/tokens`, 
      request
    );
  }

  /**
   * Get product token status and sales information
   */
  async getProductTokenStatus(productId: string): Promise<GetProductTokenStatusResponse> {
    return await this.httpClient.get<GetProductTokenStatusResponse>(
      `/products/${productId}/tokens/status`
    );
  }
}

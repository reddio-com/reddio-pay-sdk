/**
 * Product related type definitions
 */

export interface Product {
  productId: string;
  accountId: string;
  name: string;
  description?: string;
  content: string;
  active: boolean;
  productTokens: ProductToken[];
  createdAt: string;
  totalSaleCount: number;
  totalSaleAmount: number;
}

export interface ProductToken {
  productTokenId: string;
  productId: string;
  accountId: string;
  tokenId: string;
  price: string;
  recipientAddress: string;
  paymentRouterAddress: string;
  createdAt: string;
  chainId: string;
  chainName: string;
}

export interface CreateProductRequest {
  name: string;
  description?: string;
  content: string;
  tokenIdList: string[];
  price: string;
  recipientAddress: string;
}

export interface CreateProductResponse {
  message: string;
  product: Product;
}

export interface ListProductsResponse {
  message: string;
  products: Product[];
}

export interface AddProductTokenRequest {
  tokenId: string;
  price: string;
  recipientAddress: string;
}

export interface AddProductTokenResponse {
  message: string;
  productToken: ProductToken;
}

export interface GetProductTokenStatusResponse {
  message: string;
  status: ProductTokenStatus[];
}

export interface ProductTokenStatus {
  tokenName: string;
  totalSaleCount: number;
  totalSaleAmount: number;
}

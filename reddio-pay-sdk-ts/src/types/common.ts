/**
 * Common API response interfaces and base types
 */

export interface ApiResponse<T = any> {
  message: string;
  data?: T;
}

export interface PaginationOptions {
  limit?: number;
  offset?: number;
}

export interface PaginatedResponse<T> {
  message: string;
  data: T[];
  totalCount: number;
  totalPages: number;
  currentPage: number;
  pageSize: number;
}

export interface ClientConfig {
  baseURL: string;
  apiKey: string;
  timeout?: number;
}

export interface AuthResponse {
  accessToken: string;
  expiresIn?: number;
}

export class ReddioPayError extends Error {
  constructor(
    message: string,
    public statusCode?: number,
    public code?: string
  ) {
    super(message);
    this.name = 'ReddioPayError';
  }
}

export class AuthenticationError extends ReddioPayError {
  constructor(message: string = 'Authentication failed') {
    super(message, 401, 'AUTH_ERROR');
  }
}

export class NetworkError extends ReddioPayError {
  constructor(message: string = 'Network request failed') {
    super(message, 0, 'NETWORK_ERROR');
  }
}

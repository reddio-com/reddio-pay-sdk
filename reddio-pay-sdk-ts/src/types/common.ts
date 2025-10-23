/**
 * Common API response interfaces and base types
 */

/**
 * Default endpoint constants
 */
export const REDDIO_ENDPOINTS = {
  PRODUCTION: 'https://reddio-service-prod.reddio.com',
  DEVELOPMENT: 'https://reddio-service-dev.reddio.com',
} as const;

export const DEFAULT_CONFIG = {
  TIMEOUT: 30000,        // 30 seconds default timeout
  RETRY_ATTEMPTS: 3,     // Default retry attempts
  RETRY_DELAY: 1000,     // Retry delay (ms)
} as const;

export type ReddioEnvironment = 'prod' | 'dev';

/**
 * API response interface
 */
export interface ApiResponse<T = any> {
  message: string;
  data?: T;
}

/**
 * Pagination options
 */
export interface PaginationOptions {
  limit?: number;
  offset?: number;
}

/**
 * Paginated response
 */
export interface PaginatedResponse<T> {
  message: string;
  data: T[];
  totalCount: number;
  totalPages: number;
  currentPage: number;
  pageSize: number;
}

/**
 * Client config interface - baseURL is now optional
 */
export interface ClientConfig {
  apiKey: string;                          // Required: API key
  baseURL?: string;                        // Optional: custom API base URL
  environment?: ReddioEnvironment;         // Optional: environment, default 'production'
  timeout?: number;                        // Optional: request timeout, default 30000ms
}

/**
 * Internal resolved config
 */
export interface ResolvedClientConfig {
  apiKey: string;
  baseURL: string;                         
  timeout: number;                         
}

/**
 * Auth response (corresponds to Go SDK's LoginByAPIKeyResponse)
 */
export interface AuthResponse {
  message: string;
  access_token: string;
  refresh_token: string;
}

/**
 * Custom error classes
 */
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

export class ValidationError extends ReddioPayError {
  constructor(message: string = 'Validation failed') {
    super(message, 400, 'VALIDATION_ERROR');
  }
}

/**
 * Config resolve utility function
 */
export function resolveClientConfig(config: ClientConfig): ResolvedClientConfig {
  return {
    apiKey: config.apiKey,
    baseURL: resolveBaseURL(config),
    timeout: config.timeout ?? DEFAULT_CONFIG.TIMEOUT,
  };
}

/**
 * Resolve baseURL - only two environments: prod and dev
 */
function resolveBaseURL(config: ClientConfig): string {
  // 1. User explicitly specifies baseURL
  if (config.baseURL) {
    return config.baseURL;
  }

  // 2. Select by environment (only two environments)
  switch (config.environment) {
    case 'dev':
      return REDDIO_ENDPOINTS.DEVELOPMENT;
    case 'prod':
    default:
      return REDDIO_ENDPOINTS.PRODUCTION; // Default to production
  }
}
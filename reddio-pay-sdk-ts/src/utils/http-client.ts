import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { ClientConfig, AuthResponse, ReddioPayError, NetworkError, AuthenticationError } from '../types';

/**
 * Token manager handles JWT token authentication and automatic refresh
 */
class TokenManager {
  private accessToken: string | null = null;
  private refreshTimer: NodeJS.Timeout | null = null;
  private readonly apiKey: string;
  private readonly httpClient: AxiosInstance;

  constructor(apiKey: string, httpClient: AxiosInstance) {
    this.apiKey = apiKey;
    this.httpClient = httpClient;
  }

  /**
   * Authenticate using API key and get access token
   */
  async authenticate(): Promise<void> {
    try {
      const response = await this.httpClient.post('/auth/login', {
        apiKey: this.apiKey
      });

      const authData: AuthResponse = response.data;
      this.accessToken = authData.accessToken;
      
      // Start token refresh timer (refresh every hour)
      this.startTokenRefresh();
    } catch (error) {
      throw new AuthenticationError(`Failed to authenticate: ${error}`);
    }
  }

  /**
   * Get current access token
   */
  getToken(): string | null {
    return this.accessToken;
  }

  /**
   * Start automatic token refresh
   */
  private startTokenRefresh(): void {
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer);
    }

    this.refreshTimer = setInterval(async () => {
      try {
        await this.authenticate();
      } catch (error) {
        console.error('Failed to refresh token:', error);
      }
    }, 60 * 60 * 1000); // Refresh every hour
  }

  /**
   * Stop token refresh and cleanup
   */
  destroy(): void {
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer);
      this.refreshTimer = null;
    }
    this.accessToken = null;
  }
}

/**
 * HTTP client with built-in authentication and error handling
 */
export class HttpClient {
  private readonly axiosInstance: AxiosInstance;
  private readonly tokenManager: TokenManager;

  constructor(config: ClientConfig) {
    // Create axios instance
    this.axiosInstance = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout || 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Initialize token manager
    this.tokenManager = new TokenManager(config.apiKey, this.axiosInstance);

    // Setup request interceptor to add auth header
    this.axiosInstance.interceptors.request.use((requestConfig) => {
      const token = this.tokenManager.getToken();
      if (token) {
        requestConfig.headers = requestConfig.headers || {};
        requestConfig.headers.Authorization = `Bearer ${token}`;
      }
      return requestConfig;
    });

    // Setup response interceptor for error handling
    this.axiosInstance.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response) {
          // Server responded with error status
          throw new ReddioPayError(
            error.response.data?.message || error.message,
            error.response.status,
            error.response.data?.code
          );
        } else if (error.request) {
          // Request was made but no response received
          throw new NetworkError('No response received from server');
        } else {
          // Something else happened
          throw new ReddioPayError(error.message);
        }
      }
    );
  }

  /**
   * Initialize the client (authenticate)
   */
  async initialize(): Promise<void> {
    await this.tokenManager.authenticate();
  }

  /**
   * Make GET request
   */
  async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.axiosInstance.get(url, config);
    return response.data;
  }

  /**
   * Make POST request
   */
  async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.axiosInstance.post(url, data, config);
    return response.data;
  }

  /**
   * Make PUT request
   */
  async put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.axiosInstance.put(url, data, config);
    return response.data;
  }

  /**
   * Make DELETE request
   */
  async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response: AxiosResponse<T> = await this.axiosInstance.delete(url, config);
    return response.data;
  }

  /**
   * Cleanup and destroy the client
   */
  destroy(): void {
    this.tokenManager.destroy();
  }
}

import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';
import { ResolvedClientConfig, AuthResponse, AuthenticationError, NetworkError } from '../types/common';

const LOG_LEVEL = process.env.LOG_LEVEL || 'info';
function debugLog(...args: any[]) {
  if (LOG_LEVEL === 'debug') {
    console.log(...args);
  }
}

export class HttpClient {
  private axiosInstance: AxiosInstance;
  private config: ResolvedClientConfig;
  private accessToken?: string;
  private refreshTokenValue?: string;

  constructor(config: ResolvedClientConfig) {
    this.config = config;
    
    this.axiosInstance = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors(): void {
    // Request interceptor - debug log
    this.axiosInstance.interceptors.request.use(
      (config) => {
        debugLog('[HTTP REQUEST]');
        debugLog(`   Method: ${config.method?.toUpperCase()}`);
        debugLog(`   URL: ${config.baseURL}${config.url}`);
        debugLog(`   Headers: ${JSON.stringify(config.headers, null, 2)}`);
        if (config.data) {
          debugLog(`   Body: ${JSON.stringify(config.data, null, 2)}`);
        }
        if (config.params) {
          debugLog(`   Query: ${JSON.stringify(config.params, null, 2)}`);
        }

        if (this.accessToken) {
          config.headers.Authorization = `Bearer ${this.accessToken}`;
          debugLog(`   Added Authorization: Bearer ${this.accessToken.substring(0, 20)}...`);
        }
        return config;
      },
      (error) => {
        debugLog('[HTTP REQUEST ERROR]', error);
        return Promise.reject(error);
      }
    );

    // Response interceptor - debug log
    this.axiosInstance.interceptors.response.use(
      (response) => {
        debugLog('[HTTP RESPONSE]');
        debugLog(`   Status: ${response.status} ${response.statusText}`);
        debugLog(`   Headers: ${JSON.stringify(response.headers, null, 2)}`);
        debugLog(`   Body: ${JSON.stringify(response.data, null, 2)}`);
        return response;
      },
      async (error) => {
        debugLog('[HTTP RESPONSE ERROR]');
        if (error.response) {
          debugLog(`   Status: ${error.response.status} ${error.response.statusText}`);
          debugLog(`   Headers: ${JSON.stringify(error.response.headers, null, 2)}`);
          debugLog(`   Body: ${JSON.stringify(error.response.data, null, 2)}`);
        } else if (error.request) {
          debugLog(`   Network error: ${error.message}`);
          debugLog(`   Request config: ${JSON.stringify(error.config, null, 2)}`);
        } else {
          debugLog(`   Unknown error: ${error.message}`);
        }

        if (error.response?.status === 401 && this.refreshTokenValue) {
          try {
            debugLog('Trying to refresh token...');
            await this.refreshToken();
            debugLog('Retrying original request...');
            return this.axiosInstance.request(error.config);
          } catch (refreshError) {
            throw new AuthenticationError('Token refresh failed');
          }
        }
        
        if (error.response) {
          const statusCode = error.response.status;
          const message = error.response.data?.message || error.message;
          
          if (statusCode === 401) {
            throw new AuthenticationError(message);
          } else {
            throw new Error(`Request failed with status code ${statusCode}: ${message}`);
          }
        } else if (error.request) {
          throw new NetworkError('Network request failed');
        } else {
          throw new Error(error.message);
        }
      }
    );
  }

  async initialize(): Promise<void> {
    try {
      debugLog('[AUTHENTICATION] Start authentication...');
      const response = await this.axiosInstance.post('/accounts/apikeys/login', {
        api_key: this.config.apiKey
      });
      const data = response.data;
      this.accessToken = data.access_token;
      this.refreshTokenValue = data.refresh_token;
      debugLog('[AUTHENTICATION] Success');
      debugLog(`   Access Token: ${this.accessToken?.substring(0, 20)}...`);
      debugLog(`   Refresh Token: ${this.refreshTokenValue?.substring(0, 20)}...`);
      this.startTokenRefresh();
    } catch (error: unknown) {
      debugLog('[AUTHENTICATION] Failed');
      if (axios.isAxiosError(error)) {
        if (error.response) {
          const statusCode = error.response.status;
          const message = error.response.data?.message || error.message;
          throw new AuthenticationError(`Authentication failed (${statusCode}): ${message}`);
        } else {
          throw new AuthenticationError(`Authentication failed: ${error.message}`);
        }
      } else {
        const errorMessage = error instanceof Error ? error.message : 'Unknown authentication error';
        throw new AuthenticationError(`Authentication failed: ${errorMessage}`);
      }
    }
  }

  private startTokenRefresh(): void {
    debugLog('[TOKEN REFRESH] Start auto refresh (every 55 minutes)');
    setInterval(async () => {
      try {
        debugLog('[TOKEN REFRESH] Refreshing token...');
        await this.refreshToken();
      } catch (error: unknown) {
        const errorMessage = error instanceof Error ? error.message : 'Unknown refresh error';
        debugLog('[TOKEN REFRESH] Token refresh failed:', errorMessage);
      }
    }, 55 * 60 * 1000);
  }

  private async refreshToken(): Promise<void> {
    try {
      debugLog('[TOKEN REFRESH] Start refreshing token...');
      const response = await this.axiosInstance.post('/accounts/apikeys/login', {
        api_key: this.config.apiKey
      });
      const data = response.data;
      this.accessToken = data.access_token;
      this.refreshTokenValue = data.refresh_token;
      debugLog('[TOKEN REFRESH] Token refreshed');
      debugLog(`   New Access Token: ${this.accessToken?.substring(0, 20)}...`);
    } catch (error: unknown) {
      if (axios.isAxiosError(error)) {
        const message = error.response?.data?.message || error.message;
        throw new Error(`Token refresh failed: ${message}`);
      } else {
        const errorMessage = error instanceof Error ? error.message : 'Unknown refresh error';
        throw new Error(`Token refresh failed: ${errorMessage}`);
      }
    }
  }

  getBaseURL(): string {
    return this.config.baseURL;
  }

  async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    debugLog(`[API CALL] GET ${url}`);
    const response = await this.axiosInstance.get(url, config);
    return response.data;
  }

  async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    debugLog(`[API CALL] POST ${url}`);
    const response = await this.axiosInstance.post(url, data, config);
    return response.data;
  }

  async put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    debugLog(`[API CALL] PUT ${url}`);
    const response = await this.axiosInstance.put(url, data, config);
    return response.data;
  }

  async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    debugLog(`[API CALL] DELETE ${url}`);
    const response = await this.axiosInstance.delete(url, config);
    return response.data;
  }

  destroy(): void {
    debugLog('[CLEANUP] HttpClient resources cleaned up');
    this.accessToken = undefined;
    this.refreshTokenValue = undefined;
  }
}
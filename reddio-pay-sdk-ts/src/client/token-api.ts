import { HttpClient } from '../utils/http-client';
import { Token, ListTokensResponse } from '../types';

/**
 * Token API client for managing supported tokens
 */
export class TokenApi {
  constructor(private httpClient: HttpClient) {}

  /**
   * Get list of all supported tokens
   */
  async listTokens(): Promise<Token[]> {
    const response = await this.httpClient.get<ListTokensResponse>('/tokens');
    return response.tokens;
  }
}

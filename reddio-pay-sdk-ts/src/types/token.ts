/**
 * Token related type definitions
 */

export interface Token {
  tokenId: string;
  name: string;
  symbol: string;
  contractAddress: string;
  decimals: number;
  chainId: number;
  chainName: string;
  chainSymbol: string;
  explorerUrl: string;
  iconUrl: string;
  tokenType: string;
  isActive: boolean;
  currencyType: string;
  createdAt: string;
}

export interface ListTokensResponse {
  message: string;
  tokens: Token[];
  count: number;
}

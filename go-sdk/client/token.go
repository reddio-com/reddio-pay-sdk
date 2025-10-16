package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Token represents a single token information
type Token struct {
	TokenID         string `json:"token_id"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	ContractAddress string `json:"contract_address"`
	Decimals        int    `json:"decimals"`
	ChainID         int    `json:"chain_id"`
	ChainName       string `json:"chain_name"`
	ChainSymbol     string `json:"chain_symbol"`
	ExplorerURL     string `json:"explorer_url"`
	IconURL         string `json:"icon_url"`
	TokenType       string `json:"token_type"`
	IsActive        bool   `json:"is_active"`
	CurrencyType    string `json:"currency_type"`
	CreatedAt       string `json:"created_at"`
}

// ListTokensResponse represents the response for listing tokens
type ListTokensResponse struct {
	Count  int     `json:"count"`
	Tokens []*Token `json:"tokens"`
}

// ListTokens retrieves all supported tokens
func (c *Client) ListTokens() (*ListTokensResponse, error) {
	url := c.url + "/tokens"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置认证头
	if c.tokenHolder != nil {
		req.Header.Set("Authorization", "Bearer "+c.tokenHolder.getToken())
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response ListTokensResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

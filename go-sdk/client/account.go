package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AccountResponse represents the account information response
type AccountResponse struct {
	Email       string `json:"email"`
	Webhook     string `json:"webhook,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
	CompanyURL  string `json:"company_url,omitempty"`
	Activated   bool   `json:"activated"`
	CreatedAt   string `json:"created_at"`
}

// UpdateWebhookRequest represents the request body for updating webhook
type UpdateWebhookRequest struct {
	Webhook string `json:"webhook"`
}

// UpdateWebhookResponse represents the response for updating webhook
type UpdateWebhookResponse struct {
	Message string `json:"message"`
}

// UpdateAccountInfoRequest represents the request body for updating account company info
type UpdateAccountInfoRequest struct {
	CompanyName string `json:"company_name"`
	CompanyURL  string `json:"company_url"`
}

// UpdateAccountInfoResponse represents the response for updating account company info
type UpdateAccountInfoResponse struct {
	Message string `json:"message"`
}

// BalanceRequest represents the request body for balance query
type BalanceRequest struct {
	WalletAddress string `json:"wallet_address"`
	ChainID       int    `json:"chain_id"`
	TokenSymbol   string `json:"token_symbol,omitempty"` // 可选，如果不提供则查询所有支持的代币
}

// TokenBalance represents a single token balance
type TokenBalance struct {
	TokenID          string `json:"token_id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	ContractAddress  string `json:"contract_address"`
	Decimals         int    `json:"decimals"`
	Balance          string `json:"balance"`           // 原始余额（wei等单位）
	FormattedBalance string `json:"formatted_balance"` // 格式化后的余额（人类可读）
	ChainID          int    `json:"chain_id"`
	ChainName        string `json:"chain_name"`
	ChainSymbol      string `json:"chain_symbol"`
	IconURL          string `json:"icon_url"` // 代币图标URL
}

// BalanceResponse represents the response for balance query
type BalanceResponse struct {
	WalletAddress string         `json:"wallet_address"`
	ChainID       int            `json:"chain_id"`
	ChainName     string         `json:"chain_name"`
	Balances      []TokenBalance `json:"balances"`
}

// AccountAddress represents an account address information
type AccountAddress struct {
	AccountID        string `json:"account_id"`
	TokenID          string `json:"token_id"`
	RecipientAddress string `json:"recipient_address"`
	RefName          string `json:"ref_name"`
	CreatedAt        string `json:"created_at"`
}

// GetAccountInfo retrieves account information for the authenticated user
func (c *Client) GetAccountInfo() (*AccountResponse, error) {
	url := c.url + "/accounts/info"
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
	var response AccountResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// UpdateWebhook updates the webhook URL for the authenticated account
func (c *Client) UpdateWebhook(webhookURL string) (*UpdateWebhookResponse, error) {
	// 构建请求
	reqBody := UpdateWebhookRequest{
		Webhook: webhookURL,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := c.url + "/accounts/webhook"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
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

	// 解析响应
	var response UpdateWebhookResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, response.Message)
	}

	return &response, nil
}

// UpdateAccountInfo updates the company information for the authenticated account
func (c *Client) UpdateAccountInfo(companyName, companyURL string) (*UpdateAccountInfoResponse, error) {
	// 构建请求
	reqBody := UpdateAccountInfoRequest{
		CompanyName: companyName,
		CompanyURL:  companyURL,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := c.url + "/accounts/info"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
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

	// 解析响应
	var response UpdateAccountInfoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, response.Message)
	}

	return &response, nil
}

// GetTokenBalances queries token balances for a wallet address on specified chain
func (c *Client) GetTokenBalances(walletAddress string, chainID int, tokenSymbol string) (*BalanceResponse, error) {
	// 构建请求
	reqBody := BalanceRequest{
		WalletAddress: walletAddress,
		ChainID:       chainID,
		TokenSymbol:   tokenSymbol,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := c.url + "/accounts/wallet/info"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

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
	var response BalanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ListAccountAddresses retrieves all account addresses for the authenticated account
func (c *Client) ListAccountAddresses() ([]*AccountAddress, error) {
	url := c.url + "/accounts/addresses"
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
	var addresses []*AccountAddress
	if err := json.Unmarshal(body, &addresses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return addresses, nil
}

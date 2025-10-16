package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateProductRequest represents the request for creating a product
type CreateProductRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Content          string   `json:"content"`
	TokenIDList      []string `json:"token_ids"`
	Price            string   `json:"price"`
	RecipientAddress string   `json:"recipient_address"`
}

// ProductToken represents a product token information
type ProductToken struct {
	ProductTokenID       string `json:"product_token_id"`
	ProductID            string `json:"product_id"`
	AccountID            string `json:"account_id"`
	TokenID              string `json:"token_id"`
	Price                string `json:"price"`
	RecipientAddress     string `json:"recipient_address"`
	PaymentRouterAddress string `json:"payment_router_address"`
	CreatedAt            string `json:"created_at"`
	ChainID              string `json:"chain_id"`
	ChainName            string `json:"chain_name"`
}

// Product represents a product information
type Product struct {
	ProductID       string         `json:"product_id"`
	AccountID       string         `json:"account_id"`
	Name            string         `json:"name"`
	Description     string         `json:"description,omitempty"`
	Content         string         `json:"content"`
	Active          bool           `json:"active"`
	ProductTokens   []*ProductToken `json:"product_tokens"`
	CreatedAt       string         `json:"created_at"`
	TotalSaleCount  int64          `json:"total_sale_count"`
	TotalSaleAmount float64        `json:"total_sale_amount"`
}

// CreateProductResponse represents the response for creating a product
type CreateProductResponse struct {
	Message string  `json:"message"`
	Product *Product `json:"product"`
}

// ListProductsResponse represents the response for listing products
type ListProductsResponse struct {
	Message  string     `json:"message"`
	Products []*Product `json:"products"`
}

// AddProductTokenRequest represents the request for adding a token to a product
type AddProductTokenRequest struct {
	TokenID          string `json:"token_id"`
	Price            string `json:"price"`
	RecipientAddress string `json:"recipient_address"`
}

// AddProductTokenResponse represents the response for adding a token to a product
type AddProductTokenResponse struct {
	Message      string       `json:"message"`
	ProductToken *ProductToken `json:"product_token"`
}

// ProductTokenStatus represents the status of a product token
type ProductTokenStatus struct {
	ProductName     string `json:"product_name,omitempty"`
	TokenName       string `json:"token_name"`
	ChainName       string `json:"chain_name"`
	ProductTokenID  string `json:"product_token_id"`
	TotalSaleCount  int64  `json:"total_sale_count"`
	TotalSaleAmount int64  `json:"total_sale_amount"`
	CreatedAt       string `json:"created_at"`
	Decimals        int    `json:"decimals"`
	TokenID         string `json:"token_id"`
	Desc            string `json:"desc"`
}

// GetProductTokenStatusResponse represents the response for getting product token status
type GetProductTokenStatusResponse struct {
	Message string                 `json:"message"`
	Status  []*ProductTokenStatus `json:"status"`
}

// CreateProduct creates a new product
func (c *Client) CreateProduct(req *CreateProductRequest) (*CreateProductResponse, error) {
	// 构建请求
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := c.url + "/products"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	if c.tokenHolder != nil {
		httpReq.Header.Set("Authorization", "Bearer "+c.tokenHolder.getToken())
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
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
	var response CreateProductResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, response.Message)
	}

	return &response, nil
}

// ListProducts retrieves all products for the authenticated account
func (c *Client) ListProducts() (*ListProductsResponse, error) {
	url := c.url + "/products"
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
	var response ListProductsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetProduct retrieves a specific product by ID
func (c *Client) GetProduct(productID string) (*Product, error) {
	url := c.url + "/products/" + productID
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
	var product Product
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &product, nil
}

// AddProductToken adds a token to a specific product
func (c *Client) AddProductToken(productID string, req *AddProductTokenRequest) (*AddProductTokenResponse, error) {
	// 构建请求
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建 HTTP 请求
	url := c.url + "/products/" + productID + "/token"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	if c.tokenHolder != nil {
		httpReq.Header.Set("Authorization", "Bearer "+c.tokenHolder.getToken())
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
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
	var response AddProductTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, response.Message)
	}

	return &response, nil
}

// GetProductTokenStatus retrieves the status of all tokens for a specific product
func (c *Client) GetProductTokenStatus(productID string) (*GetProductTokenStatusResponse, error) {
	url := c.url + "/products/" + productID + "/token/status"
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
	var response GetProductTokenStatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

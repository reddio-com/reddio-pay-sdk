package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Payment represents a payment information
type Payment struct {
	PaymentID        string `json:"payment_id"`
	AccountID        string `json:"account_id"`
	TokenID          string `json:"token_id"`
	ProductID        string `json:"product_id"`
	ProductTokenID   string `json:"product_token_id"`
	Count            int    `json:"count"`
	Status           string `json:"status"`
	PayerEmail       string `json:"payer_email,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	PaidAt           string `json:"paid_at,omitempty"`
	ClosedAt         string `json:"closed_at,omitempty"`
	CloseReason      string `json:"close_reason,omitempty"`
	TransactionHash  string `json:"transaction_hash,omitempty"`
	BlockNumber      int64  `json:"block_number,omitempty"`
	GasUsed          int64  `json:"gas_used,omitempty"`
	GasPrice         string `json:"gas_price,omitempty"`
	TotalAmount      string `json:"total_amount"`
	FeeAmount        string `json:"fee_amount"`
	RecipientAmount  string `json:"recipient_amount"`
}

// ListPaymentsResponse represents the response for listing payments
type ListPaymentsResponse struct {
	Message  string     `json:"message"`
	Payments []*Payment `json:"payments"`
}

// ListPaymentsResponseWithPagination represents the response for listing payments with pagination
type ListPaymentsResponseWithPagination struct {
	Message     string     `json:"message"`
	Payments    []*Payment `json:"payments"`
	TotalCount  int        `json:"total_count"`
	TotalPages  int        `json:"total_pages"`
	CurrentPage int        `json:"current_page"`
	PageSize    int        `json:"page_size"`
}

// ExternalSendNotifyForPaymentSuccessRequest represents the request for external payment success notification
type ExternalSendNotifyForPaymentSuccessRequest struct {
	PaymentID string `json:"payment_id"`
	Email     string `json:"email"`
}

// SendNotifyForPaymentSuccessResponse represents the response for sending payment success notification
type SendNotifyForPaymentSuccessResponse struct {
	Message string `json:"message"`
}

// ListPaymentsByAccountAndProductID retrieves all payments for the authenticated account and specific product
func (c *Client) ListPaymentsByAccountAndProductID(productID string) (*ListPaymentsResponse, error) {
	url := c.url + "/payments/product/" + productID
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
	var response ListPaymentsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ListPaymentsByAccount retrieves all payments for the authenticated account with pagination
func (c *Client) ListPaymentsByAccount(limit, offset int) (*ListPaymentsResponseWithPagination, error) {
	// 构建 URL 和查询参数
	baseURL := c.url + "/payments/list"
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// 添加查询参数
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}
	if offset >= 0 {
		params.Add("offset", strconv.Itoa(offset))
	}
	u.RawQuery = params.Encode()

	// 创建请求
	req, err := http.NewRequest("GET", u.String(), nil)
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
	var response ListPaymentsResponseWithPagination
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetPaymentByID retrieves a specific payment by ID
func (c *Client) GetPaymentByID(paymentID string) (*Payment, error) {
	url := c.url + "/payments/" + paymentID
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
	var payment Payment
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &payment, nil
}

// ExternalSendNotifyForPaymentSuccess sets up email notification for payment success
func (c *Client) ExternalSendNotifyForPaymentSuccess(req *ExternalSendNotifyForPaymentSuccessRequest) (*SendNotifyForPaymentSuccessResponse, error) {
	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建请求
	url := c.url + "/external/payments/success/notify"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
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

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response SendNotifyForPaymentSuccessResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// PaymentReceiver represents a payment receiver information
type PaymentReceiver struct {
	Type             string `json:"type"` // "fee" or "merchant"
	RecipientAddress string `json:"recipient_address"`
	Amount           string `json:"amount"` // wei format
	Rate             string `json:"rate"`   // percentage
}


// ExternalCreatePaymentRequest represents the request for creating an external payment
type ExternalCreatePaymentRequest struct {
	ProductID      string `json:"product_id"`
	ProductTokenID string `json:"product_token_id"`
	Count          int    `json:"count"`
}

// ExternalCreatePaymentResponse represents the response for creating an external payment
type ExternalCreatePaymentResponse struct {
	Message          string             `json:"message"`
	PaymentID        string             `json:"payment_id"`
	PayLink          string             `json:"pay_link"`
	ContractAddress  string             `json:"contract_address"`
	PaymentReceivers []*PaymentReceiver `json:"payment_receivers"`
	TokenAddress     string             `json:"token_address"`
	Decimals         int                `json:"decimals"`
}


// ExternalCreatePayment creates a new external payment
func (c *Client) ExternalCreatePayment(req *ExternalCreatePaymentRequest) (*ExternalCreatePaymentResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	url := c.url + "/external/payments"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.tokenHolder != nil {
		httpReq.Header.Set("Authorization", "Bearer "+c.tokenHolder.getToken())
	}
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	var response ExternalCreatePaymentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &response, nil
}

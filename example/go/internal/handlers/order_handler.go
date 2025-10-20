package handlers

import (
	"encoding/json"
	"net/http"
	"order-system/internal/services"
	"strconv"

	"github.com/gorilla/mux"
)

// OrderHandler order handler
type OrderHandler struct {
	orderService *services.OrderService
}

// NewOrderHandler creates order handler
func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// writeJSONResponse writes JSON response
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}



// CreateOrder creates order API
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerName   string `json:"customer_name"`
		CustomerEmail  string `json:"customer_email"`
		ProductID      string `json:"product_id"`
		ProductTokenID string `json:"product_token_id"`
		Quantity       int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.CustomerName == "" || req.CustomerEmail == "" || req.ProductID == "" || req.ProductTokenID == "" || req.Quantity < 1 {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Missing required fields or invalid quantity",
		})
		return
	}

	order, err := h.orderService.CreateOrder(
		req.CustomerName,
		req.CustomerEmail,
		req.ProductID,
		req.ProductTokenID,
		req.Quantity,
	)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"order": order,
	})
}

// GetOrder gets order details API
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid order ID",
		})
		return
	}

	order, err := h.orderService.GetOrder(orderID)
	if err != nil {
		writeJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "Order not found",
		})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"order": order,
	})
}

// GetOrderStatus gets order status API
func (h *OrderHandler) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid order ID",
		})
		return
	}

	order, err := h.orderService.GetOrder(orderID)
	if err != nil {
		writeJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": "Order not found",
		})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"order_number":     order.OrderNumber,
		"status":           order.Status,
		"reddio_status":    order.ReddioStatus,
		"reddio_pay_link":  order.ReddioPayLink,
		"transaction_hash": order.TransactionHash,
		"paid_at":          order.PaidAt,
	})
}

// CheckPayment checks payment status API
func (h *OrderHandler) CheckPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid order ID",
		})
		return
	}

	err = h.orderService.CheckPaymentStatus(orderID)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	order, err := h.orderService.GetOrder(orderID)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get updated order",
		})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message":          "Payment status updated",
		"order_number":     order.OrderNumber,
		"status":           order.Status,
		"reddio_status":    order.ReddioStatus,
		"transaction_hash": order.TransactionHash,
		"paid_at":          order.PaidAt,
	})
}

// ListOrders gets order list API
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	status := r.URL.Query().Get("status")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)
	if pageInt < 1 {
		pageInt = 1
	}
	if limitInt < 1 || limitInt > 100 {
		limitInt = 10
	}

	orders, total, err := h.orderService.ListOrders(pageInt, limitInt, status)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get orders",
		})
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
		"pagination": map[string]interface{}{
			"page":        pageInt,
			"limit":       limitInt,
			"total":       total,
			"total_pages": (total + int64(limitInt) - 1) / int64(limitInt),
		},
	})
}
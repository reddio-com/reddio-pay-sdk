package main

import (
	"log"
	"net/http"

	"order-system/internal/config"
	"order-system/internal/database"
	"order-system/internal/handlers"
	"order-system/internal/services"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Reddio Pay SDK client
	client, err := services.NewReddioPayClient(cfg)
	if err != nil {
		log.Fatal("Failed to initialize Reddio Pay client:", err)
	}

	// Initialize services
	orderService := services.NewOrderService(db, cfg)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	// Initialize SDK demo service
	sdkDemoService := services.NewSDKDemoService(client)

	// Setup routes
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	api.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	api.HandleFunc("/orders/{id}/status", orderHandler.GetOrderStatus).Methods("GET")
	api.HandleFunc("/orders/{id}/check-payment", orderHandler.CheckPayment).Methods("POST")
	api.HandleFunc("/orders", orderHandler.ListOrders).Methods("GET")
	// Demonstrate SDK APIs on startup
	log.Println("Starting Reddio Pay Go SDK API demonstration...")
	if err := sdkDemoService.DemoAllAPIs(); err != nil {
		log.Printf("Error during SDK demonstration: %v", err)
		log.Println("Continuing to start server...")
	}

	log.Println("Order system starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

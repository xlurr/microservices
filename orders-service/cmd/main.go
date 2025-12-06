package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"orders-service/internal/handlers"
	"orders-service/internal/repository"
	_ "orders-service/docs"
)

// @title Orders Service API
// @version 1.0
// @description Orders microservice
// @host localhost:8082
// @basePath /api

func main() {
	router := mux.NewRouter()

	// Инициализируем JSON репозиторий
	repo, err := repository.NewJSONOrderRepository("./data/orders.json")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	handler := handlers.NewOrderHandler(repo)

	// API endpoints
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", handler.GetAllOrders).Methods("GET")
	api.HandleFunc("/orders/{id}", handler.GetOrderByID).Methods("GET")
	api.HandleFunc("/orders/user/{userId}", handler.GetOrdersByUserID).Methods("GET")
	api.HandleFunc("/orders/{id}", handler.UpdateOrder).Methods("PUT")
	api.HandleFunc("/orders/{id}", handler.DeleteOrder).Methods("DELETE")
	api.HandleFunc("/orders/user/{userId}", handler.DeleteOrdersByUserID).Methods("DELETE")

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Starting orders-service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

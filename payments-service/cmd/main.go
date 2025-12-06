package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"payments-service/internal/handlers"
	"payments-service/internal/repository"
	_ "payments-service/docs"
)

// @title Payments Service API
// @version 1.0
// @description Payments microservice
// @host localhost:8083
// @basePath /api

func main() {
	router := mux.NewRouter()

	// Инициализируем JSON репозиторий
	repo, err := repository.NewJSONPaymentRepository("./data/payments.json")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	handler := handlers.NewPaymentHandler(repo)

	// API endpoints
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/payments", handler.CreatePayment).Methods("POST")
	api.HandleFunc("/payments", handler.GetAllPayments).Methods("GET")
	api.HandleFunc("/payments/{id}", handler.GetPaymentByID).Methods("GET")
	api.HandleFunc("/payments/user/{userId}", handler.GetPaymentsByUserID).Methods("GET")
	api.HandleFunc("/payments/{id}", handler.UpdatePayment).Methods("PUT")
	api.HandleFunc("/payments/user/{userId}", handler.DeletePaymentsByUserID).Methods("DELETE")

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8083/swagger/doc.json"),
	))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Starting payments-service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

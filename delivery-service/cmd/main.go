package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"delivery-service/internal/handlers"
	"delivery-service/internal/repository"
	_ "delivery-service/docs"
)

// @title Delivery Service API
// @version 1.0
// @description Delivery microservice
// @host localhost:8084
// @basePath /api

func main() {
	router := mux.NewRouter()

	// Инициализируем JSON репозиторий
	repo, err := repository.NewJSONDeliveryRepository("./data/deliveries.json")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	handler := handlers.NewDeliveryHandler(repo)

	// API endpoints
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/deliveries", handler.CreateDelivery).Methods("POST")
	api.HandleFunc("/deliveries", handler.GetAllDeliveries).Methods("GET")
	api.HandleFunc("/deliveries/{id}", handler.GetDeliveryByID).Methods("GET")
	api.HandleFunc("/deliveries/user/{userId}", handler.GetDeliveriesByUserID).Methods("GET")
	api.HandleFunc("/deliveries/{id}", handler.UpdateDelivery).Methods("PUT")
	api.HandleFunc("/deliveries/user/{userId}", handler.DeleteDeliveriesByUserID).Methods("DELETE")

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8084/swagger/doc.json"),
	))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("Starting delivery-service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

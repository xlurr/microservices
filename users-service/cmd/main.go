package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"users-service/internal/handlers"
	"users-service/internal/repository"
	_ "users-service/docs"
)

// @title Users Service API
// @version 1.0
// @description Users microservice
// @host localhost:8081
// @basePath /api

func main() {
	router := mux.NewRouter()

	// Инициализируем JSON репозиторий
	repo, err := repository.NewJSONUserRepository("./data/users.json")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	handler := handlers.NewUserHandler(repo)

	// API endpoints
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/users", handler.CreateUser).Methods("POST")
	api.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", handler.GetUserByID).Methods("GET")
	api.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")
	api.HandleFunc("/users/{id}/exists", handler.UserExists).Methods("GET")

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
	))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting users-service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

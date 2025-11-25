package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "users-service/docs" // ✅ Импортируйте docs пакет
	"users-service/internal/handlers"
	"users-service/internal/repository"
)

// @title Users Microservice API
// @version 1.0
// @description API для управления пользователями
// @host localhost:8081
// @BasePath /
// @schemes http https
func main() {
	router := mux.NewRouter()

	repo := repository.NewInMemoryUserRepository()
	handler := handlers.NewUserHandler(repo)

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", handler.CreateUser).Methods("POST")
	api.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", handler.GetUserByID).Methods("GET")
	api.HandleFunc("/users/{id}", handler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", handler.DeleteUser).Methods("DELETE")

	// ✅ Правильная конфигурация Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Users service starting on :8081")
	log.Println("Swagger UI: http://localhost:8081/swagger/")
	log.Fatal(srv.ListenAndServe())
}

package main

import (
"log"
"net/http"
"time"

"github.com/gorilla/mux"
httpSwagger "github.com/swaggo/http-swagger"

"delivery-service/internal/handlers"
"delivery-service/internal/repository"
_ "delivery-service/docs"
)

// @title Delivery Microservice API
// @version 1.0
// @description API для управления доставками
// @host localhost:8084
// @BasePath /
// @schemes http https
func main() {
router := mux.NewRouter()

repo := repository.NewInMemoryDeliveryRepository()
handler := handlers.NewDeliveryHandler(repo)

api := router.PathPrefix("/api").Subrouter()
api.HandleFunc("/deliveries", handler.CreateDelivery).Methods("POST")
api.HandleFunc("/deliveries", handler.GetAllDeliveries).Methods("GET")

router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
httpSwagger.URL("http://localhost:8084/swagger/doc.json"),
httpSwagger.DeepLinking(true),
httpSwagger.DocExpansion("list"),
))

router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
w.Write([]byte("OK"))
}).Methods("GET")

srv := &http.Server{
Handler:      router,
Addr:         ":8084",
WriteTimeout: 15 * time.Second,
ReadTimeout:  15 * time.Second,
}

log.Println("Delivery service starting on :8084")
log.Println("Swagger: http://localhost:8084/swagger/")
log.Fatal(srv.ListenAndServe())
}

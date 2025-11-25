package main

import (
"log"
"net/http"
"time"

"github.com/gorilla/mux"
httpSwagger "github.com/swaggo/http-swagger"

"orders-service/internal/handlers"
"orders-service/internal/repository"
_ "orders-service/docs"
)

// @title Orders Microservice API
// @version 1.0
// @description API для управления заказами
// @host localhost:8082
// @BasePath /
// @schemes http https
func main() {
router := mux.NewRouter()

repo := repository.NewInMemoryOrderRepository()
handler := handlers.NewOrderHandler(repo)

api := router.PathPrefix("/api").Subrouter()
api.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
api.HandleFunc("/orders", handler.GetAllOrders).Methods("GET")
api.HandleFunc("/orders/{id}", handler.GetOrderByID).Methods("GET")
api.HandleFunc("/orders/{id}/status", handler.UpdateOrderStatus).Methods("PUT")

router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
httpSwagger.DeepLinking(true),
httpSwagger.DocExpansion("list"),
))

router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
w.Write([]byte("OK"))
}).Methods("GET")

srv := &http.Server{
Handler:      router,
Addr:         ":8082",
WriteTimeout: 15 * time.Second,
ReadTimeout:  15 * time.Second,
}

log.Println("Orders service starting on :8082")
log.Println("Swagger: http://localhost:8082/swagger/")
log.Fatal(srv.ListenAndServe())
}

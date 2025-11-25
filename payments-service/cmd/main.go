package main

import (
"log"
"net/http"
"time"

"github.com/gorilla/mux"
httpSwagger "github.com/swaggo/http-swagger"

"payments-service/internal/handlers"
"payments-service/internal/repository"
_ "payments-service/docs"
)

// @title Payments Microservice API
// @version 1.0
// @description API для управления платежами
// @host localhost:8083
// @BasePath /
// @schemes http https
func main() {
router := mux.NewRouter()

repo := repository.NewInMemoryPaymentRepository()
handler := handlers.NewPaymentHandler(repo)

api := router.PathPrefix("/api").Subrouter()
api.HandleFunc("/payments", handler.CreatePayment).Methods("POST")
api.HandleFunc("/payments", handler.GetAllPayments).Methods("GET")

router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
httpSwagger.URL("http://localhost:8083/swagger/doc.json"),
httpSwagger.DeepLinking(true),
httpSwagger.DocExpansion("list"),
))

router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
w.Write([]byte("OK"))
}).Methods("GET")

srv := &http.Server{
Handler:      router,
Addr:         ":8083",
WriteTimeout: 15 * time.Second,
ReadTimeout:  15 * time.Second,
}

log.Println("Payments service starting on :8083")
log.Println("Swagger: http://localhost:8083/swagger/")
log.Fatal(srv.ListenAndServe())
}

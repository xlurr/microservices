#!/bin/bash

#############################################
# Скрипт для автоматической генерации
# микросервисной архитектуры на Go
# с Docker и Swagger интеграцией
#############################################

set -e  # Остановка при ошибке
set -u  # Остановка при использовании неинициализированных переменных

# Цвета для вывода
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  Генерация структуры микросервисов на Go${NC}"
echo -e "${BLUE}================================================${NC}"

# Проверка наличия Go
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}⚠️  Go не установлен. Установите Go 1.20+${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Go версия: $(go version)${NC}"

# Создание основной структуры директорий
echo -e "\n${BLUE}📁 Создание структуры директорий...${NC}"

# Функция для создания директорий с индикацией прогресса
create_structure() {
    local service=$1
    echo -e "${GREEN}  → Создание структуры для ${service}${NC}"
    
    mkdir -p "${service}/cmd"
    mkdir -p "${service}/internal/handlers"
    mkdir -p "${service}/internal/models"
    mkdir -p "${service}/internal/repository"
    
    if [[ "$service" != "users-service" ]]; then
        mkdir -p "${service}/internal/client"
    fi
}

# Создание структуры для всех сервисов
for service in users-service orders-service payments-service delivery-service; do
    create_structure "$service"
done

echo -e "${GREEN}✓ Структура директорий создана${NC}"

#############################################
# USERS SERVICE
#############################################

echo -e "\n${BLUE}📝 Генерация Users Service...${NC}"

# users-service/internal/models/user.go
cat > users-service/internal/models/user.go << 'EOF'
package models

import "time"

type User struct {
    ID        int64     `json:"id"`
    FirstName string    `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string    `json:"last_name" validate:"required,min=2,max=50"`
    Email     string    `json:"email" validate:"required,email"`
    Age       int       `json:"age" validate:"required,gte=18,lte=120"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
    FirstName string `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string `json:"last_name" validate:"required,min=2,max=50"`
    Email     string `json:"email" validate:"required,email"`
    Age       int    `json:"age" validate:"required,gte=18,lte=120"`
}

type ErrorResponse struct {
    Status  int          `json:"status"`
    Message string       `json:"message"`
    Errors  []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
EOF

# users-service/internal/repository/user_repository.go
cat > users-service/internal/repository/user_repository.go << 'EOF'
package repository

import (
    "errors"
    "sync"
    "users-service/internal/models"
)

type UserRepository interface {
    Create(user *models.User) error
    GetAll() ([]*models.User, error)
    GetByID(id int64) (*models.User, error)
    Update(user *models.User) error
    Delete(id int64) error
}

type InMemoryUserRepository struct {
    mu      sync.RWMutex
    users   map[int64]*models.User
    nextID  int64
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users:  make(map[int64]*models.User),
        nextID: 1,
    }
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    user.ID = r.nextID
    r.users[user.ID] = user
    r.nextID++
    
    return nil
}

func (r *InMemoryUserRepository) GetAll() ([]*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    users := make([]*models.User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }
    
    return users, nil
}

func (r *InMemoryUserRepository) GetByID(id int64) (*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    user, exists := r.users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    
    return user, nil
}

func (r *InMemoryUserRepository) Update(user *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.users[user.ID]; !exists {
        return errors.New("user not found")
    }
    
    r.users[user.ID] = user
    return nil
}

func (r *InMemoryUserRepository) Delete(id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.users[id]; !exists {
        return errors.New("user not found")
    }
    
    delete(r.users, id)
    return nil
}
EOF

# users-service/internal/handlers/user_handler.go
cat > users-service/internal/handlers/user_handler.go << 'EOF'
package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/go-playground/validator/v10"
    
    "users-service/internal/models"
    "users-service/internal/repository"
)

type UserHandler struct {
    repo     repository.UserRepository
    validate *validator.Validate
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
    return &UserHandler{
        repo:     repo,
        validate: validator.New(),
    }
}

// CreateUser godoc
// @Summary Создать нового пользователя
// @Description Создание пользователя с валидацией данных
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body", nil)
        return
    }
    
    if err := h.validate.Struct(req); err != nil {
        validationErrors := translateValidationErrors(err.(validator.ValidationErrors))
        respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    user := &models.User{
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Email:     req.Email,
        Age:       req.Age,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := h.repo.Create(user); err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create user", nil)
        return
    }
    
    respondWithJSON(w, http.StatusCreated, user)
}

// GetAllUsers godoc
// @Summary Получить всех пользователей
// @Description Список всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.repo.GetAll()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to fetch users", nil)
        return
    }
    
    respondWithJSON(w, http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Description Информация о пользователе
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID", nil)
        return
    }
    
    user, err := h.repo.GetByID(id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "User not found", nil)
        return
    }
    
    respondWithJSON(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Обновить пользователя
// @Description Обновление данных пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.CreateUserRequest true "Обновленные данные"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID", nil)
        return
    }
    
    var req models.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body", nil)
        return
    }
    
    if err := h.validate.Struct(req); err != nil {
        validationErrors := translateValidationErrors(err.(validator.ValidationErrors))
        respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    user, err := h.repo.GetByID(id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "User not found", nil)
        return
    }
    
    user.FirstName = req.FirstName
    user.LastName = req.LastName
    user.Email = req.Email
    user.Age = req.Age
    user.UpdatedAt = time.Now()
    
    if err := h.repo.Update(user); err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to update user", nil)
        return
    }
    
    respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Description Удаление пользователя по ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID", nil)
        return
    }
    
    if err := h.repo.Delete(id); err != nil {
        respondWithError(w, http.StatusNotFound, "User not found", nil)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string, fieldErrors []models.FieldError) {
    errResponse := models.ErrorResponse{
        Status:  code,
        Message: message,
        Errors:  fieldErrors,
    }
    respondWithJSON(w, code, errResponse)
}

func translateValidationErrors(errs validator.ValidationErrors) []models.FieldError {
    var fieldErrors []models.FieldError
    for _, err := range errs {
        fieldErrors = append(fieldErrors, models.FieldError{
            Field:   err.Field(),
            Message: getErrorMessage(err),
        })
    }
    return fieldErrors
}

func getErrorMessage(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return "This field is required"
    case "email":
        return "Invalid email format"
    case "min":
        return fmt.Sprintf("Minimum length is %s", err.Param())
    case "max":
        return fmt.Sprintf("Maximum length is %s", err.Param())
    case "gte":
        return fmt.Sprintf("Must be >= %s", err.Param())
    case "lte":
        return fmt.Sprintf("Must be <= %s", err.Param())
    default:
        return "Invalid value"
    }
}
EOF

# users-service/cmd/main.go
cat > users-service/cmd/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    httpSwagger "github.com/swaggo/http-swagger"
    
    "users-service/internal/handlers"
    "users-service/internal/repository"
)

// @title Users Microservice API
// @version 1.0
// @description API для управления пользователями
// @host localhost:8081
// @BasePath /
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
    
    // Swagger UI
    router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
        httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
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
EOF

# users-service/go.mod
cat > users-service/go.mod << 'EOF'
module users-service

go 1.21

require (
    github.com/gorilla/mux v1.8.1
    github.com/go-playground/validator/v10 v10.16.0
    github.com/swaggo/http-swagger v1.3.4
    github.com/swaggo/swag v1.16.2
)
EOF

# users-service/Dockerfile
cat > users-service/Dockerfile << 'EOF'
# Стадия сборки
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Установка swag для генерации документации
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-w -s' \
    -o /app/users-service ./cmd/main.go

# Финальная стадия
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/users-service /users-service

USER appuser

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --quiet --tries=1 --spider http://localhost:8081/health || exit 1

ENTRYPOINT ["/users-service"]
EOF

echo -e "${GREEN}✓ Users Service создан${NC}"

#############################################
# ORDERS SERVICE
#############################################

echo -e "\n${BLUE}📝 Генерация Orders Service...${NC}"

# orders-service/internal/models/order.go
cat > orders-service/internal/models/order.go << 'EOF'
package models

import "time"

type Order struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id" validate:"required"`
    Items       []string  `json:"items" validate:"required,min=1"`
    TotalAmount float64   `json:"total_amount"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
    UserID int64    `json:"user_id" validate:"required"`
    Items  []string `json:"items" validate:"required,min=1"`
}

type UpdateStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=created payment_pending paid delivered cancelled"`
}

type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}
EOF

# orders-service/internal/repository/order_repository.go
cat > orders-service/internal/repository/order_repository.go << 'EOF'
package repository

import (
    "errors"
    "sync"
    "orders-service/internal/models"
)

type OrderRepository interface {
    Create(order *models.Order) error
    GetAll() ([]*models.Order, error)
    GetByID(id int64) (*models.Order, error)
    UpdateStatus(id int64, status string) error
}

type InMemoryOrderRepository struct {
    mu     sync.RWMutex
    orders map[int64]*models.Order
    nextID int64
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
    return &InMemoryOrderRepository{
        orders: make(map[int64]*models.Order),
        nextID: 1,
    }
}

func (r *InMemoryOrderRepository) Create(order *models.Order) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    order.ID = r.nextID
    r.orders[order.ID] = order
    r.nextID++
    
    return nil
}

func (r *InMemoryOrderRepository) GetAll() ([]*models.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    orders := make([]*models.Order, 0, len(r.orders))
    for _, order := range r.orders {
        orders = append(orders, order)
    }
    
    return orders, nil
}

func (r *InMemoryOrderRepository) GetByID(id int64) (*models.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    order, exists := r.orders[id]
    if !exists {
        return nil, errors.New("order not found")
    }
    
    return order, nil
}

func (r *InMemoryOrderRepository) UpdateStatus(id int64, status string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    order, exists := r.orders[id]
    if !exists {
        return errors.New("order not found")
    }
    
    order.Status = status
    return nil
}
EOF

# orders-service/internal/client/service_client.go
cat > orders-service/internal/client/service_client.go << 'EOF'
package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type ServiceClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewServiceClient(baseURL string) *ServiceClient {
    return &ServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *ServiceClient) CreatePayment(orderID int64, amount float64) error {
    payload := map[string]interface{}{
        "order_id": orderID,
        "amount":   amount,
        "status":   "pending",
    }
    
    body, _ := json.Marshal(payload)
    
    req, err := http.NewRequest("POST", 
        fmt.Sprintf("%s/api/payments", c.baseURL), 
        bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("failed to create payment: status %d", resp.StatusCode)
    }
    
    return nil
}
EOF

# orders-service/internal/handlers/order_handler.go
cat > orders-service/internal/handlers/order_handler.go << 'EOF'
package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/go-playground/validator/v10"
    
    "orders-service/internal/models"
    "orders-service/internal/repository"
    "orders-service/internal/client"
)

type OrderHandler struct {
    repo     repository.OrderRepository
    validate *validator.Validate
}

func NewOrderHandler(repo repository.OrderRepository) *OrderHandler {
    return &OrderHandler{
        repo:     repo,
        validate: validator.New(),
    }
}

// CreateOrder godoc
// @Summary Создать заказ
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Данные заказа"
// @Success 201 {object} models.Order
// @Router /api/orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var req models.CreateOrderRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    if err := h.validate.Struct(req); err != nil {
        respondWithError(w, http.StatusUnprocessableEntity, "Validation failed")
        return
    }
    
    order := &models.Order{
        UserID:      req.UserID,
        Items:       req.Items,
        TotalAmount: float64(len(req.Items)) * 100.0,
        Status:      "created",
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    if err := h.repo.Create(order); err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create order")
        return
    }
    
    // Асинхронное создание платежа
    go func() {
        paymentsURL := os.Getenv("PAYMENTS_SERVICE_URL")
        if paymentsURL == "" {
            paymentsURL = "http://payments-service:8083"
        }
        
        paymentClient := client.NewServiceClient(paymentsURL)
        if err := paymentClient.CreatePayment(order.ID, order.TotalAmount); err != nil {
            log.Printf("Failed to create payment: %v", err)
            h.repo.UpdateStatus(order.ID, "payment_failed")
        } else {
            h.repo.UpdateStatus(order.ID, "payment_pending")
        }
    }()
    
    respondWithJSON(w, http.StatusCreated, order)
}

// GetAllOrders godoc
// @Summary Получить все заказы
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Router /api/orders [get]
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
    orders, err := h.repo.GetAll()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to fetch orders")
        return
    }
    
    respondWithJSON(w, http.StatusOK, orders)
}

// GetOrderByID godoc
// @Summary Получить заказ по ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Router /api/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.ParseInt(vars["id"], 10, 64)
    
    order, err := h.repo.GetByID(id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Order not found")
        return
    }
    
    respondWithJSON(w, http.StatusOK, order)
}

// UpdateOrderStatus godoc
// @Summary Обновить статус заказа
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body models.UpdateStatusRequest true "Новый статус"
// @Success 200 {object} models.Order
// @Router /api/orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.ParseInt(vars["id"], 10, 64)
    
    var req models.UpdateStatusRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }
    
    if err := h.repo.UpdateStatus(id, req.Status); err != nil {
        respondWithError(w, http.StatusNotFound, "Order not found")
        return
    }
    
    order, _ := h.repo.GetByID(id)
    respondWithJSON(w, http.StatusOK, order)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, models.ErrorResponse{
        Status:  code,
        Message: message,
    })
}
EOF

# orders-service/cmd/main.go
cat > orders-service/cmd/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
    httpSwagger "github.com/swaggo/http-swagger"
    
    "orders-service/internal/handlers"
    "orders-service/internal/repository"
)

// @title Orders Microservice API
// @version 1.0
// @description API для управления заказами
// @host localhost:8082
// @BasePath /
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
    log.Fatal(srv.ListenAndServe())
}
EOF

# orders-service/go.mod
cat > orders-service/go.mod << 'EOF'
module orders-service

go 1.21

require (
    github.com/gorilla/mux v1.8.1
    github.com/go-playground/validator/v10 v10.16.0
    github.com/swaggo/http-swagger v1.3.4
    github.com/swaggo/swag v1.16.2
)
EOF

# orders-service/Dockerfile
cat > orders-service/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates
RUN adduser -D -g '' appuser

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 go build -ldflags='-w -s' -o /app/orders-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/orders-service /orders-service

USER appuser
EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --quiet --tries=1 --spider http://localhost:8082/health || exit 1

ENTRYPOINT ["/orders-service"]
EOF

echo -e "${GREEN}✓ Orders Service создан${NC}"

#############################################
# PAYMENTS SERVICE (упрощенная версия)
#############################################

echo -e "\n${BLUE}📝 Генерация Payments Service...${NC}"

cat > payments-service/internal/models/payment.go << 'EOF'
package models

import "time"

type Payment struct {
    ID        int64     `json:"id"`
    OrderID   int64     `json:"order_id"`
    Amount    float64   `json:"amount"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

type CreatePaymentRequest struct {
    OrderID int64   `json:"order_id" validate:"required"`
    Amount  float64 `json:"amount" validate:"required,gt=0"`
}
EOF

cat > payments-service/internal/repository/payment_repository.go << 'EOF'
package repository

import (
    "errors"
    "sync"
    "payments-service/internal/models"
)

type PaymentRepository interface {
    Create(payment *models.Payment) error
    GetAll() ([]*models.Payment, error)
    GetByID(id int64) (*models.Payment, error)
}

type InMemoryPaymentRepository struct {
    mu       sync.RWMutex
    payments map[int64]*models.Payment
    nextID   int64
}

func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
    return &InMemoryPaymentRepository{
        payments: make(map[int64]*models.Payment),
        nextID:   1,
    }
}

func (r *InMemoryPaymentRepository) Create(payment *models.Payment) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    payment.ID = r.nextID
    r.payments[payment.ID] = payment
    r.nextID++
    
    return nil
}

func (r *InMemoryPaymentRepository) GetAll() ([]*models.Payment, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    payments := make([]*models.Payment, 0, len(r.payments))
    for _, payment := range r.payments {
        payments = append(payments, payment)
    }
    
    return payments, nil
}

func (r *InMemoryPaymentRepository) GetByID(id int64) (*models.Payment, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    payment, exists := r.payments[id]
    if !exists {
        return nil, errors.New("payment not found")
    }
    
    return payment, nil
}
EOF

cat > payments-service/internal/handlers/payment_handler.go << 'EOF'
package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/go-playground/validator/v10"
    
    "payments-service/internal/models"
    "payments-service/internal/repository"
)

type PaymentHandler struct {
    repo     repository.PaymentRepository
    validate *validator.Validate
}

func NewPaymentHandler(repo repository.PaymentRepository) *PaymentHandler {
    return &PaymentHandler{
        repo:     repo,
        validate: validator.New(),
    }
}

// CreatePayment godoc
// @Summary Создать платеж
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body models.CreatePaymentRequest true "Данные платежа"
// @Success 201 {object} models.Payment
// @Router /api/payments [post]
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
    var req models.CreatePaymentRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    
    if err := h.validate.Struct(req); err != nil {
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }
    
    payment := &models.Payment{
        OrderID:   req.OrderID,
        Amount:    req.Amount,
        Status:    "pending",
        CreatedAt: time.Now(),
    }
    
    h.repo.Create(payment)
    
    json.NewEncoder(w).Encode(payment)
}

// GetAllPayments godoc
// @Summary Получить все платежи
// @Tags payments
// @Produce json
// @Success 200 {array} models.Payment
// @Router /api/payments [get]
func (h *PaymentHandler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
    payments, _ := h.repo.GetAll()
    json.NewEncoder(w).Encode(payments)
}
EOF

cat > payments-service/cmd/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    
    "github.com/gorilla/mux"
    httpSwagger "github.com/swaggo/http-swagger"
    
    "payments-service/internal/handlers"
    "payments-service/internal/repository"
)

// @title Payments Microservice API
// @version 1.0
// @host localhost:8083
// @BasePath /
func main() {
    router := mux.NewRouter()
    
    repo := repository.NewInMemoryPaymentRepository()
    handler := handlers.NewPaymentHandler(repo)
    
    api := router.PathPrefix("/api").Subrouter()
    api.HandleFunc("/payments", handler.CreatePayment).Methods("POST")
    api.HandleFunc("/payments", handler.GetAllPayments).Methods("GET")
    
    router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
        httpSwagger.URL("http://localhost:8083/swagger/doc.json"),
    ))
    
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }).Methods("GET")
    
    log.Println("Payments service starting on :8083")
    log.Fatal(http.ListenAndServe(":8083", router))
}
EOF

cat > payments-service/go.mod << 'EOF'
module payments-service

go 1.21

require (
    github.com/gorilla/mux v1.8.1
    github.com/go-playground/validator/v10 v10.16.0
    github.com/swaggo/http-swagger v1.3.4
    github.com/swaggo/swag v1.16.2
)
EOF

cat > payments-service/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go
RUN CGO_ENABLED=0 go build -o /app/payments-service ./cmd/main.go

FROM alpine:latest
COPY --from=builder /app/payments-service /payments-service
EXPOSE 8083
ENTRYPOINT ["/payments-service"]
EOF

echo -e "${GREEN}✓ Payments Service создан${NC}"

#############################################
# DELIVERY SERVICE (упрощенная версия)
#############################################

echo -e "\n${BLUE}📝 Генерация Delivery Service...${NC}"

cat > delivery-service/internal/models/delivery.go << 'EOF'
package models

import "time"

type Delivery struct {
    ID        int64     `json:"id"`
    OrderID   int64     `json:"order_id"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
EOF

cat > delivery-service/internal/repository/delivery_repository.go << 'EOF'
package repository

import (
    "sync"
    "delivery-service/internal/models"
)

type DeliveryRepository interface {
    Create(delivery *models.Delivery) error
    GetAll() ([]*models.Delivery, error)
}

type InMemoryDeliveryRepository struct {
    mu         sync.RWMutex
    deliveries map[int64]*models.Delivery
    nextID     int64
}

func NewInMemoryDeliveryRepository() *InMemoryDeliveryRepository {
    return &InMemoryDeliveryRepository{
        deliveries: make(map[int64]*models.Delivery),
        nextID:     1,
    }
}

func (r *InMemoryDeliveryRepository) Create(delivery *models.Delivery) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    delivery.ID = r.nextID
    r.deliveries[delivery.ID] = delivery
    r.nextID++
    
    return nil
}

func (r *InMemoryDeliveryRepository) GetAll() ([]*models.Delivery, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    deliveries := make([]*models.Delivery, 0, len(r.deliveries))
    for _, delivery := range r.deliveries {
        deliveries = append(deliveries, delivery)
    }
    
    return deliveries, nil
}
EOF

cat > delivery-service/internal/handlers/delivery_handler.go << 'EOF'
package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    
    "delivery-service/internal/models"
    "delivery-service/internal/repository"
)

type DeliveryHandler struct {
    repo repository.DeliveryRepository
}

func NewDeliveryHandler(repo repository.DeliveryRepository) *DeliveryHandler {
    return &DeliveryHandler{repo: repo}
}

// CreateDelivery godoc
// @Summary Создать доставку
// @Tags delivery
// @Accept json
// @Produce json
// @Success 201 {object} models.Delivery
// @Router /api/deliveries [post]
func (h *DeliveryHandler) CreateDelivery(w http.ResponseWriter, r *http.Request) {
    var req struct {
        OrderID int64 `json:"order_id"`
    }
    
    json.NewDecoder(r.Body).Decode(&req)
    
    delivery := &models.Delivery{
        OrderID:   req.OrderID,
        Status:    "pending",
        CreatedAt: time.Now(),
    }
    
    h.repo.Create(delivery)
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(delivery)
}

// GetAllDeliveries godoc
// @Summary Получить все доставки
// @Tags delivery
// @Produce json
// @Success 200 {array} models.Delivery
// @Router /api/deliveries [get]
func (h *DeliveryHandler) GetAllDeliveries(w http.ResponseWriter, r *http.Request) {
    deliveries, _ := h.repo.GetAll()
    json.NewEncoder(w).Encode(deliveries)
}
EOF

cat > delivery-service/cmd/main.go << 'EOF'
package main

import (
    "log"
    "net/http"
    
    "github.com/gorilla/mux"
    httpSwagger "github.com/swaggo/http-swagger"
    
    "delivery-service/internal/handlers"
    "delivery-service/internal/repository"
)

// @title Delivery Microservice API
// @version 1.0
// @host localhost:8084
// @BasePath /
func main() {
    router := mux.NewRouter()
    
    repo := repository.NewInMemoryDeliveryRepository()
    handler := handlers.NewDeliveryHandler(repo)
    
    api := router.PathPrefix("/api").Subrouter()
    api.HandleFunc("/deliveries", handler.CreateDelivery).Methods("POST")
    api.HandleFunc("/deliveries", handler.GetAllDeliveries).Methods("GET")
    
    router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
        httpSwagger.URL("http://localhost:8084/swagger/doc.json"),
    ))
    
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }).Methods("GET")
    
    log.Println("Delivery service starting on :8084")
    log.Fatal(http.ListenAndServe(":8084", router))
}
EOF

cat > delivery-service/go.mod << 'EOF'
module delivery-service

go 1.21

require (
    github.com/gorilla/mux v1.8.1
    github.com/swaggo/http-swagger v1.3.4
    github.com/swaggo/swag v1.16.2
)
EOF

cat > delivery-service/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go
RUN CGO_ENABLED=0 go build -o /app/delivery-service ./cmd/main.go

FROM alpine:latest
COPY --from=builder /app/delivery-service /delivery-service
EXPOSE 8084
ENTRYPOINT ["/delivery-service"]
EOF

echo -e "${GREEN}✓ Delivery Service создан${NC}"

#############################################
# DOCKER COMPOSE
#############################################

echo -e "\n${BLUE}📝 Создание docker-compose.yml...${NC}"

cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  users-service:
    build:
      context: ./users-service
      dockerfile: Dockerfile
    container_name: users-service
    ports:
      - "8081:8081"
    environment:
      - SERVICE_NAME=users-service
      - SERVICE_PORT=8081
    networks:
      - microservices-network
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8081/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    restart: unless-stopped

  orders-service:
    build:
      context: ./orders-service
      dockerfile: Dockerfile
    container_name: orders-service
    ports:
      - "8082:8082"
    environment:
      - SERVICE_NAME=orders-service
      - SERVICE_PORT=8082
      - PAYMENTS_SERVICE_URL=http://payments-service:8083
    networks:
      - microservices-network
    depends_on:
      users-service:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8082/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    restart: unless-stopped

  payments-service:
    build:
      context: ./payments-service
      dockerfile: Dockerfile
    container_name: payments-service
    ports:
      - "8083:8083"
    environment:
      - SERVICE_NAME=payments-service
      - SERVICE_PORT=8083
    networks:
      - microservices-network
    depends_on:
      orders-service:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8083/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    restart: unless-stopped

  delivery-service:
    build:
      context: ./delivery-service
      dockerfile: Dockerfile
    container_name: delivery-service
    ports:
      - "8084:8084"
    environment:
      - SERVICE_NAME=delivery-service
      - SERVICE_PORT=8084
    networks:
      - microservices-network
    depends_on:
      payments-service:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8084/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    restart: unless-stopped

networks:
  microservices-network:
    driver: bridge

volumes:
  users-data:
  orders-data:
  payments-data:
  delivery-data:
EOF

echo -e "${GREEN}✓ docker-compose.yml создан${NC}"

#############################################
# README
#############################################

echo -e "\n${BLUE}📝 Создание README.md...${NC}"

cat > README.md << 'EOF'


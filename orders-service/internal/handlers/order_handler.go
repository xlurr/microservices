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

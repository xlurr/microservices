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

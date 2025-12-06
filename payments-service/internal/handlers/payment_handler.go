package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"payments-service/internal/repository"
)

// PaymentHandler обработчик платежей
type PaymentHandler struct {
	repo *repository.JSONPaymentRepository
}

// CreatePaymentRequest структура для создания платежа
type CreatePaymentRequest struct {
	UserID  int64   `json:"userId" binding:"required"`
	OrderID int64   `json:"orderId" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,min=0"`
}

// UpdatePaymentRequest структура для обновления платежа
type UpdatePaymentRequest struct {
	Status string `json:"status" binding:"required"`
}

// NewPaymentHandler создаёт новый обработчик
func NewPaymentHandler(repo *repository.JSONPaymentRepository) *PaymentHandler {
	return &PaymentHandler{
		repo: repo,
	}
}

// CreatePayment создаёт платёж
// @Summary Создать платёж
// @Description Создать новый платёж
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body CreatePaymentRequest true "Данные платежа"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request"
// @Router /payments [post]
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	payment, err := h.repo.CreatePayment(req.UserID, req.OrderID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

// GetAllPayments получает все платежи
// @Summary Получить все платежи
// @Description Получить список всех платежей
// @Tags payments
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /payments [get]
func (h *PaymentHandler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

// GetPaymentByID получает платёж по ID
// @Summary Получить платёж по ID
// @Description Получить конкретный платёж
// @Tags payments
// @Produce json
// @Param id path int64 true "ID платежа"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Payment not found"
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	payment, err := h.repo.GetPaymentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

// GetPaymentsByUserID получает платежи пользователя
// @Summary Получить платежи пользователя
// @Description Получить все платежи конкретного пользователя
// @Tags payments
// @Produce json
// @Param userId path int64 true "ID пользователя"
// @Success 200 {array} map[string]interface{}
// @Router /payments/user/{userId} [get]
func (h *PaymentHandler) GetPaymentsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	payments, err := h.repo.GetPaymentsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

// UpdatePayment обновляет платёж
// @Summary Обновить платёж
// @Description Обновить статус платежа
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int64 true "ID платежа"
// @Param payment body UpdatePaymentRequest true "Новый статус"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Payment not found"
// @Router /payments/{id} [put]
func (h *PaymentHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	payment, err := h.repo.UpdatePaymentStatus(id, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

// DeletePaymentsByUserID удаляет платежи пользователя
// @Summary Удалить платежи пользователя
// @Description Удалить все платежи пользователя
// @Tags payments
// @Param userId path int64 true "ID пользователя"
// @Success 204
// @Router /payments/user/{userId} [delete]
func (h *PaymentHandler) DeletePaymentsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	h.repo.DeletePaymentsByUserID(userID)
	w.WriteHeader(http.StatusNoContent)
}

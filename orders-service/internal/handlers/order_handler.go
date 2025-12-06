package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"orders-service/internal/repository"
)

// OrderHandler обработчик заказов
type OrderHandler struct {
	repo *repository.JSONOrderRepository
}

// CreateOrderRequest структура для создания заказа
type CreateOrderRequest struct {
	UserID      int64    `json:"userId" binding:"required"`
	Items       []string `json:"items" binding:"required"`
	TotalAmount float64  `json:"totalAmount" binding:"required,min=0"`
}

// UpdateOrderRequest структура для обновления заказа
type UpdateOrderRequest struct {
	Status string `json:"status" binding:"required"`
}

// NewOrderHandler создаёт новый обработчик
func NewOrderHandler(repo *repository.JSONOrderRepository) *OrderHandler {
	return &OrderHandler{
		repo: repo,
	}
}

// CreateOrder создаёт заказ
// @Summary Создать заказ
// @Description Создать новый заказ
// @Tags orders
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Данные заказа"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	order, err := h.repo.CreateOrder(req.UserID, req.Items, req.TotalAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetAllOrders получает все заказы
// @Summary Получить все заказы
// @Description Получить список всех заказов
// @Tags orders
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /orders [get]
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

// GetOrderByID получает заказ по ID
// @Summary Получить заказ по ID
// @Description Получить конкретный заказ
// @Tags orders
// @Produce json
// @Param id path int64 true "ID заказа"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Order not found"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	order, err := h.repo.GetOrderByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetOrdersByUserID получает заказы пользователя
// @Summary Получить заказы пользователя
// @Description Получить все заказы конкретного пользователя
// @Tags orders
// @Produce json
// @Param userId path int64 true "ID пользователя"
// @Success 200 {array} map[string]interface{}
// @Router /orders/user/{userId} [get]
func (h *OrderHandler) GetOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	orders, err := h.repo.GetOrdersByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrder обновляет заказ
// @Summary Обновить заказ
// @Description Обновить статус заказа
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int64 true "ID заказа"
// @Param order body UpdateOrderRequest true "Новый статус"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Order not found"
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	order, err := h.repo.UpdateOrderStatus(id, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// DeleteOrder удаляет заказ
// @Summary Удалить заказ
// @Description Удалить заказ по ID
// @Tags orders
// @Param id path int64 true "ID заказа"
// @Success 204
// @Failure 404 {string} string "Order not found"
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteOrdersByUserID удаляет заказы пользователя
// @Summary Удалить заказы пользователя
// @Description Удалить все заказы пользователя
// @Tags orders
// @Param userId path int64 true "ID пользователя"
// @Success 204
// @Router /orders/user/{userId} [delete]
func (h *OrderHandler) DeleteOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	h.repo.DeleteOrdersByUserID(userID)
	w.WriteHeader(http.StatusNoContent)
}

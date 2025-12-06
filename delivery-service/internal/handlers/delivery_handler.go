package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"delivery-service/internal/repository"
)

// DeliveryHandler обработчик доставок
type DeliveryHandler struct {
	repo *repository.JSONDeliveryRepository
}

// CreateDeliveryRequest структура для создания доставки
type CreateDeliveryRequest struct {
	UserID     int64  `json:"userId" binding:"required"`
	OrderID    int64  `json:"orderId" binding:"required"`
	Address    string `json:"address" binding:"required"`
	TrackingID string `json:"trackingId" binding:"required"`
}

// UpdateDeliveryRequest структура для обновления доставки
type UpdateDeliveryRequest struct {
	Status string `json:"status" binding:"required"`
}

// NewDeliveryHandler создаёт новый обработчик
func NewDeliveryHandler(repo *repository.JSONDeliveryRepository) *DeliveryHandler {
	return &DeliveryHandler{
		repo: repo,
	}
}

// CreateDelivery создаёт доставку
// @Summary Создать доставку
// @Description Создать новую доставку
// @Tags deliveries
// @Accept json
// @Produce json
// @Param delivery body CreateDeliveryRequest true "Данные доставки"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request"
// @Router /deliveries [post]
func (h *DeliveryHandler) CreateDelivery(w http.ResponseWriter, r *http.Request) {
	var req CreateDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	delivery, err := h.repo.CreateDelivery(req.UserID, req.OrderID, req.Address, req.TrackingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(delivery)
}

// GetAllDeliveries получает все доставки
// @Summary Получить все доставки
// @Description Получить список всех доставок
// @Tags deliveries
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /deliveries [get]
func (h *DeliveryHandler) GetAllDeliveries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

// GetDeliveryByID получает доставку по ID
// @Summary Получить доставку по ID
// @Description Получить конкретную доставку
// @Tags deliveries
// @Produce json
// @Param id path int64 true "ID доставки"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Delivery not found"
// @Router /deliveries/{id} [get]
func (h *DeliveryHandler) GetDeliveryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	delivery, err := h.repo.GetDeliveryByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delivery)
}

// GetDeliveriesByUserID получает доставки пользователя
// @Summary Получить доставки пользователя
// @Description Получить все доставки конкретного пользователя
// @Tags deliveries
// @Produce json
// @Param userId path int64 true "ID пользователя"
// @Success 200 {array} map[string]interface{}
// @Router /deliveries/user/{userId} [get]
func (h *DeliveryHandler) GetDeliveriesByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	deliveries, err := h.repo.GetDeliveriesByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deliveries)
}

// UpdateDelivery обновляет доставку
// @Summary Обновить доставку
// @Description Обновить статус доставки
// @Tags deliveries
// @Accept json
// @Produce json
// @Param id path int64 true "ID доставки"
// @Param delivery body UpdateDeliveryRequest true "Новый статус"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Delivery not found"
// @Router /deliveries/{id} [put]
func (h *DeliveryHandler) UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	delivery, err := h.repo.UpdateDeliveryStatus(id, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delivery)
}

// DeleteDeliveriesByUserID удаляет доставки пользователя
// @Summary Удалить доставки пользователя
// @Description Удалить все доставки пользователя
// @Tags deliveries
// @Param userId path int64 true "ID пользователя"
// @Success 204
// @Router /deliveries/user/{userId} [delete]
func (h *DeliveryHandler) DeleteDeliveriesByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	h.repo.DeleteDeliveriesByUserID(userID)
	w.WriteHeader(http.StatusNoContent)
}

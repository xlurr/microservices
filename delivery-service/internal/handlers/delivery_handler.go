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

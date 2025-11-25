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

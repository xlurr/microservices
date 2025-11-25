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

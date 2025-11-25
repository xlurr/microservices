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

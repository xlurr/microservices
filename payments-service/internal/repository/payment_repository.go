package repository

import (
	"fmt"
	"sync"
	"time"

	"payments-service/internal/models"
	"payments-service/internal/utils"
)

// JSONPaymentRepository - репозиторий с JSON-хранилищем
type JSONPaymentRepository struct {
	mu       sync.RWMutex
	storage  *utils.FileStorage
	payments map[int64]*models.Payment
	nextID   int64
	filePath string
}

// NewJSONPaymentRepository создаёт новый репозиторий
func NewJSONPaymentRepository(filePath string) (*JSONPaymentRepository, error) {
	repo := &JSONPaymentRepository{
		storage:  utils.NewFileStorage(filePath),
		payments: make(map[int64]*models.Payment),
		nextID:   1,
		filePath: filePath,
	}

	if err := repo.storage.EnsureFile(); err != nil {
		return nil, err
	}

	if err := repo.LoadFromFile(); err != nil {
		return nil, err
	}

	return repo, nil
}

// LoadFromFile загружает платежи из JSON
func (r *JSONPaymentRepository) LoadFromFile() error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	var payments []*models.Payment
	if err := r.storage.LoadJSON(&payments); err != nil {
		return err
	}

	r.payments = make(map[int64]*models.Payment)
	r.nextID = 1

	for _, payment := range payments {
		r.payments[payment.ID] = payment
		if payment.ID >= r.nextID {
			r.nextID = payment.ID + 1
		}
	}

	return nil
}

// SaveToFile сохраняет платежи в JSON
func (r *JSONPaymentRepository) SaveToFile() error {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	payments := make([]*models.Payment, 0, len(r.payments))
	for _, payment := range r.payments {
		payments = append(payments, payment)
	}

	return r.storage.SaveJSON(payments)
}

// CreatePayment создаёт новый платёж (с проверкой пользователя)
func (r *JSONPaymentRepository) CreatePayment(userID, orderID int64, amount float64) (*models.Payment, error) {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	payment := &models.Payment{
		ID:        r.nextID,
		UserID:    userID,
		OrderID:   orderID,
		Amount:    amount,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r.payments[payment.ID] = payment
	r.nextID++

	if err := r.SaveToFile(); err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPaymentByID получает платёж по ID
func (r *JSONPaymentRepository) GetPaymentByID(id int64) (*models.Payment, error) {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	payment, exists := r.payments[id]
	if !exists {
		return nil, fmt.Errorf("payment not found")
	}

	return payment, nil
}

// GetPaymentsByUserID получает все платежи пользователя
func (r *JSONPaymentRepository) GetPaymentsByUserID(userID int64) ([]*models.Payment, error) {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	var userPayments []*models.Payment
	for _, payment := range r.payments {
		if payment.UserID == userID {
			userPayments = append(userPayments, payment)
		}
	}

	return userPayments, nil
}

// UpdatePaymentStatus обновляет статус платежа
func (r *JSONPaymentRepository) UpdatePaymentStatus(id int64, status string) (*models.Payment, error) {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	payment, exists := r.payments[id]
	if !exists {
		return nil, fmt.Errorf("payment not found")
	}

	payment.Status = status
	payment.UpdatedAt = time.Now()

	if err := r.SaveToFile(); err != nil {
		return nil, err
	}

	return payment, nil
}

// DeletePaymentsByUserID удаляет все платежи пользователя
func (r *JSONPaymentRepository) DeletePaymentsByUserID(userID int64) error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	paymentsToDelete := []int64{}
	for id, payment := range r.payments {
		if payment.UserID == userID {
			paymentsToDelete = append(paymentsToDelete, id)
		}
	}

	for _, id := range paymentsToDelete {
		delete(r.payments, id)
	}

	if len(paymentsToDelete) > 0 {
		return r.SaveToFile()
	}

	return nil
}

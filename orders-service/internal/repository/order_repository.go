package repository

import (
	"fmt"
	"sync"
	"time"

	"orders-service/internal/models"
	"orders-service/internal/utils"
)

// JSONOrderRepository - репозиторий с JSON-хранилищем
type JSONOrderRepository struct {
	mu       sync.RWMutex
	storage  *utils.FileStorage
	orders   map[int64]*models.Order
	nextID   int64
	filePath string
}

// NewJSONOrderRepository создаёт новый репозиторий
func NewJSONOrderRepository(filePath string) (*JSONOrderRepository, error) {
	repo := &JSONOrderRepository{
		storage:  utils.NewFileStorage(filePath),
		orders:   make(map[int64]*models.Order),
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

// LoadFromFile загружает заказы из JSON
func (r *JSONOrderRepository) LoadFromFile() error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	var orders []*models.Order
	if err := r.storage.LoadJSON(&orders); err != nil {
		return err
	}

	r.orders = make(map[int64]*models.Order)
	r.nextID = 1

	for _, order := range orders {
		r.orders[order.ID] = order
		if order.ID >= r.nextID {
			r.nextID = order.ID + 1
		}
	}

	return nil
}

// SaveToFile сохраняет заказы в JSON
func (r *JSONOrderRepository) SaveToFile() error {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	orders := make([]*models.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}

	return r.storage.SaveJSON(orders)
}

// CreateOrder создаёт новый заказ (с проверкой пользователя через HTTP)
func (r *JSONOrderRepository) CreateOrder(userID int64, items []string, totalAmount float64) (*models.Order, error) {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	order := &models.Order{
		ID:          r.nextID,
		UserID:      userID,
		Items:       items,
		TotalAmount: totalAmount,
		Status:      "created",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	r.orders[order.ID] = order
	r.nextID++

	if err := r.SaveToFile(); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderByID получает заказ по ID
func (r *JSONOrderRepository) GetOrderByID(id int64) (*models.Order, error) {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

// GetOrdersByUserID получает все заказы пользователя
func (r *JSONOrderRepository) GetOrdersByUserID(userID int64) ([]*models.Order, error) {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	var userOrders []*models.Order
	for _, order := range r.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}

	return userOrders, nil
}

// DeleteOrdersByUserID удаляет все заказы пользователя
func (r *JSONOrderRepository) DeleteOrdersByUserID(userID int64) error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	ordersToDelete := []int64{}
	for id, order := range r.orders {
		if order.UserID == userID {
			ordersToDelete = append(ordersToDelete, id)
		}
	}

	for _, id := range ordersToDelete {
		delete(r.orders, id)
	}

	if len(ordersToDelete) > 0 {
		return r.SaveToFile()
	}

	return nil
}

// DeleteOrder удаляет один заказ
func (r *JSONOrderRepository) DeleteOrder(id int64) error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	if _, exists := r.orders[id]; !exists {
		return fmt.Errorf("order not found")
	}

	delete(r.orders, id)

	return r.SaveToFile()
}

// UpdateOrderStatus обновляет статус заказа
func (r *JSONOrderRepository) UpdateOrderStatus(id int64, status string) (*models.Order, error) {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}

	order.Status = status
	order.UpdatedAt = time.Now()

	if err := r.SaveToFile(); err != nil {
		return nil, err
	}

	return order, nil
}

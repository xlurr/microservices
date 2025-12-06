package repository

import (
	"fmt"
	"sync"
	"time"

	"delivery-service/internal/models"
	"delivery-service/internal/utils"
)

// JSONDeliveryRepository - репозиторий с JSON-хранилищем
type JSONDeliveryRepository struct {
	mu         sync.RWMutex
	storage    *utils.FileStorage
	deliveries map[int64]*models.Delivery
	nextID     int64
	filePath   string
}

// NewJSONDeliveryRepository создаёт новый репозиторий
func NewJSONDeliveryRepository(filePath string) (*JSONDeliveryRepository, error) {
	repo := &JSONDeliveryRepository{
		storage:    utils.NewFileStorage(filePath),
		deliveries: make(map[int64]*models.Delivery),
		nextID:     1,
		filePath:   filePath,
	}

	if err := repo.storage.EnsureFile(); err != nil {
		return nil, err
	}

	if err := repo.LoadFromFile(); err != nil {
		return nil, err
	}

	return repo, nil
}

// LoadFromFile загружает доставки из JSON
func (r *JSONDeliveryRepository) LoadFromFile() error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	var deliveries []*models.Delivery
	if err := r.storage.LoadJSON(&deliveries); err != nil {
		return err
	}

	r.deliveries = make(map[int64]*models.Delivery)
	r.nextID = 1

	for _, delivery := range deliveries {
		r.deliveries[delivery.ID] = delivery
		if delivery.ID >= r.nextID {
			r.nextID = delivery.ID + 1
		}
	}

	return nil
}

// SaveToFile сохраняет доставки в JSON
func (r *JSONDeliveryRepository) SaveToFile() error {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	deliveries := make([]*models.Delivery, 0, len(r.deliveries))
	for _, delivery := range r.deliveries {
		deliveries = append(deliveries, delivery)
	}

	return r.storage.SaveJSON(deliveries)
}

// CreateDelivery создаёт новую доставку (с проверкой пользователя)
func (r *JSONDeliveryRepository) CreateDelivery(userID, orderID int64, address, trackingID string) (*models.Delivery, error) {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	delivery := &models.Delivery{
		ID:         r.nextID,
		UserID:     userID,
		OrderID:    orderID,
		Address:    address,
		Status:     "pending",
		TrackingID: trackingID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	r.deliveries[delivery.ID] = delivery
	r.nextID++

	if err := r.SaveToFile(); err != nil {
		return nil, err
	}

	return delivery, nil
}

// GetDeliveryByID получает доставку по ID
func (r *JSONDeliveryRepository) GetDeliveryByID(id int64) (*models.Delivery, error) {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	delivery, exists := r.deliveries[id]
	if !exists {
		return nil, fmt.Errorf("delivery not found")
	}

	return delivery, nil
}

// GetDeliveriesByUserID получает все доставки пользователя
func (r *JSONDeliveryRepository) GetDeliveriesByUserID(userID int64) ([]*models.Delivery, error) {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	var userDeliveries []*models.Delivery
	for _, delivery := range r.deliveries {
		if delivery.UserID == userID {
			userDeliveries = append(userDeliveries, delivery)
		}
	}

	return userDeliveries, nil
}

// UpdateDeliveryStatus обновляет статус доставки
func (r *JSONDeliveryRepository) UpdateDeliveryStatus(id int64, status string) (*models.Delivery, error) {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	delivery, exists := r.deliveries[id]
	if !exists {
		return nil, fmt.Errorf("delivery not found")
	}

	delivery.Status = status
	delivery.UpdatedAt = time.Now()

	if err := r.SaveToFile(); err != nil {
		return nil, err
	}

	return delivery, nil
}

// DeleteDeliveriesByUserID удаляет все доставки пользователя
func (r *JSONDeliveryRepository) DeleteDeliveriesByUserID(userID int64) error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	deliveriesToDelete := []int64{}
	for id, delivery := range r.deliveries {
		if delivery.UserID == userID {
			deliveriesToDelete = append(deliveriesToDelete, id)
		}
	}

	for _, id := range deliveriesToDelete {
		delete(r.deliveries, id)
	}

	if len(deliveriesToDelete) > 0 {
		return r.SaveToFile()
	}

	return nil
}

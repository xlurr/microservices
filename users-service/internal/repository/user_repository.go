package repository

import (
    "errors"
    "sync"
    "users-service/internal/models"
)

type UserRepository interface {
    Create(user *models.User) error
    GetAll() ([]*models.User, error)
    GetByID(id int64) (*models.User, error)
    Update(user *models.User) error
    Delete(id int64) error
}

type InMemoryUserRepository struct {
    mu      sync.RWMutex
    users   map[int64]*models.User
    nextID  int64
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users:  make(map[int64]*models.User),
        nextID: 1,
    }
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    user.ID = r.nextID
    r.users[user.ID] = user
    r.nextID++
    
    return nil
}

func (r *InMemoryUserRepository) GetAll() ([]*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    users := make([]*models.User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }
    
    return users, nil
}

func (r *InMemoryUserRepository) GetByID(id int64) (*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    user, exists := r.users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    
    return user, nil
}

func (r *InMemoryUserRepository) Update(user *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.users[user.ID]; !exists {
        return errors.New("user not found")
    }
    
    r.users[user.ID] = user
    return nil
}

func (r *InMemoryUserRepository) Delete(id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if _, exists := r.users[id]; !exists {
        return errors.New("user not found")
    }
    
    delete(r.users, id)
    return nil
}

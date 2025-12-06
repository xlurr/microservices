package repository

import (
	"fmt"
	"sync"
	"time"

	"users-service/internal/models"
	"users-service/internal/utils"
)

// JSONUserRepository - репозиторий с JSON-хранилищем
type JSONUserRepository struct {
	mu       sync.RWMutex
	storage  *utils.FileStorage
	users    map[int64]*models.User
	nextID   int64
	filePath string
}

// NewJSONUserRepository создаёт новый репозиторий с JSON-хранилищем
func NewJSONUserRepository(filePath string) (*JSONUserRepository, error) {
	repo := &JSONUserRepository{
		storage:  utils.NewFileStorage(filePath),
		users:    make(map[int64]*models.User),
		nextID:   1,
		filePath: filePath,
	}

	// Создаём файл если его нет
	if err := repo.storage.EnsureFile(); err != nil {
		return nil, err
	}

	// Загружаем данные из файла
	if err := repo.LoadFromFile(); err != nil {
		return nil, err
	}

	return repo, nil
}

// LoadFromFile загружает пользователей из JSON файла
func (r *JSONUserRepository) LoadFromFile() error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	var users []*models.User
	if err := r.storage.LoadJSON(&users); err != nil {
		return err
	}

	r.users = make(map[int64]*models.User)
	r.nextID = 1

	for _, user := range users {
		r.users[user.ID] = user
		if user.ID >= r.nextID {
			r.nextID = user.ID + 1
		}
	}

	return nil
}

// SaveToFile сохраняет пользователей в JSON файл
func (r *JSONUserRepository) SaveToFile() error {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return r.storage.SaveJSON(users)
}

// CreateUser создаёт нового пользователя
func (r *JSONUserRepository) CreateUser(email string, name string, age int) (*models.User, error) {
	// Валидация @ в email
	if len(email) == 0 || email[0] == '@' || email[len(email)-1] == '@' {
		return nil, fmt.Errorf("invalid email: must contain @ character in valid position")
	}

	atCount := 0
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atCount++
			atIndex = i
		}
	}

	if atCount != 1 {
		return nil, fmt.Errorf("invalid email: must contain exactly one @ character")
	}

	if atIndex == 0 || atIndex == len(email)-1 {
		return nil, fmt.Errorf("invalid email: @ must not be at the beginning or end")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем уникальность email
	for _, u := range r.users {
		if u.Email == email {
			return nil, fmt.Errorf("email already exists")
		}
	}

	user := &models.User{
		ID:        r.nextID,
		Email:     email,
		Name:      name,
		Age:       age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r.users[user.ID] = user
	r.nextID++

	// Сохраняем в файл
	if err := r.SaveToFile(); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID получает пользователя по ID
func (r *JSONUserRepository) GetUserByID(id int64) (*models.User, error) {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetAllUsers получает всех пользователей
func (r *JSONUserRepository) GetAllUsers() ([]*models.User, error) {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser обновляет пользователя
func (r *JSONUserRepository) UpdateUser(id int64, email, name string, age int) (*models.User, error) {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Если email менялся - валидируем
	if email != user.Email {
		atCount := 0
		for _, c := range email {
			if c == '@' {
				atCount++
			}
		}
		if atCount != 1 {
			return nil, fmt.Errorf("invalid email: must contain exactly one @ character")
		}
	}

	user.Email = email
	user.Name = name
	user.Age = age
	user.UpdatedAt = time.Now()

	// Сохраняем в файл
	if err := r.SaveToFile(); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser удаляет пользователя
func (r *JSONUserRepository) DeleteUser(id int64) error {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(r.users, id)

	// Сохраняем в файл
	return r.SaveToFile()
}

// UserExists проверяет существование пользователя
func (r *JSONUserRepository) UserExists(id int64) bool {
	//r.mu.RLock()
	//defer r.mu.RUnlock()

	_, exists := r.users[id]
	return exists
}

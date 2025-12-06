package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"users-service/internal/repository"
)

// UserHandler обработчик пользователей
type UserHandler struct {
	repo *repository.JSONUserRepository
}

// CreateUserRequest структура для создания пользователя
type CreateUserRequest struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"required,min=0,max=150"`
}

// UpdateUserRequest структура для обновления пользователя
type UpdateUserRequest struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Age   int    `json:"age" binding:"required,min=0,max=150"`
}

// NewUserHandler создаёт новый обработчик
func NewUserHandler(repo *repository.JSONUserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// CreateUser создаёт пользователя
// @Summary Создать пользователя
// @Description Создать нового пользователя с валидацией @
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "Данные пользователя"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid email"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.repo.CreateUser(req.Email, req.Name, req.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetAllUsers получает всех пользователей
// @Summary Получить всех пользователей
// @Description Получить список всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID получает пользователя по ID
// @Summary Получить пользователя по ID
// @Description Получить конкретного пользователя
// @Tags users
// @Produce json
// @Param id path int64 true "ID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser обновляет пользователя
// @Summary Обновить пользователя
// @Description Обновить данные пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int64 true "ID пользователя"
// @Param user body UpdateUserRequest true "Новые данные"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.repo.UpdateUser(id, req.Email, req.Name, req.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser удаляет пользователя
// @Summary Удалить пользователя
// @Description Удалить пользователя (каскадно удаляет заказы, платежи, доставки)
// @Tags users
// @Param id path int64 true "ID пользователя"
// @Success 204
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UserExists проверяет существование пользователя
// @Summary Проверить существование пользователя
// @Description Проверить, существует ли пользователь
// @Tags users
// @Param id path int64 true "ID пользователя"
// @Success 200 {object} map[string]bool
// @Router /users/{id}/exists [get]
func (h *UserHandler) UserExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	exists := h.repo.UserExists(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
}

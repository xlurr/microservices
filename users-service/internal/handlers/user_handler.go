package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/go-playground/validator/v10"
    
    "users-service/internal/models"
    "users-service/internal/repository"
)

type UserHandler struct {
    repo     repository.UserRepository
    validate *validator.Validate
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
    return &UserHandler{
        repo:     repo,
        validate: validator.New(),
    }
}

// CreateUser godoc
// @Summary Создать нового пользователя
// @Description Создание пользователя с валидацией данных
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body", nil)
        return
    }
    
    if err := h.validate.Struct(req); err != nil {
        validationErrors := translateValidationErrors(err.(validator.ValidationErrors))
        respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    user := &models.User{
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Email:     req.Email,
        Age:       req.Age,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := h.repo.Create(user); err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to create user", nil)
        return
    }
    
    respondWithJSON(w, http.StatusCreated, user)
}

// GetAllUsers godoc
// @Summary Получить всех пользователей
// @Description Список всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.repo.GetAll()
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to fetch users", nil)
        return
    }
    
    respondWithJSON(w, http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Description Информация о пользователе
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID", nil)
        return
    }
    
    user, err := h.repo.GetByID(id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "User not found", nil)
        return
    }
    
    respondWithJSON(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Обновить пользователя
// @Description Обновление данных пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.CreateUserRequest true "Обновленные данные"
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID", nil)
        return
    }
    
    var req models.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request body", nil)
        return
    }
    
    if err := h.validate.Struct(req); err != nil {
        validationErrors := translateValidationErrors(err.(validator.ValidationErrors))
        respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    user, err := h.repo.GetByID(id)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "User not found", nil)
        return
    }
    
    user.FirstName = req.FirstName
    user.LastName = req.LastName
    user.Email = req.Email
    user.Age = req.Age
    user.UpdatedAt = time.Now()
    
    if err := h.repo.Update(user); err != nil {
        respondWithError(w, http.StatusInternalServerError, "Failed to update user", nil)
        return
    }
    
    respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Description Удаление пользователя по ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseInt(vars["id"], 10, 64)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid user ID", nil)
        return
    }
    
    if err := h.repo.Delete(id); err != nil {
        respondWithError(w, http.StatusNotFound, "User not found", nil)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string, fieldErrors []models.FieldError) {
    errResponse := models.ErrorResponse{
        Status:  code,
        Message: message,
        Errors:  fieldErrors,
    }
    respondWithJSON(w, code, errResponse)
}

func translateValidationErrors(errs validator.ValidationErrors) []models.FieldError {
    var fieldErrors []models.FieldError
    for _, err := range errs {
        fieldErrors = append(fieldErrors, models.FieldError{
            Field:   err.Field(),
            Message: getErrorMessage(err),
        })
    }
    return fieldErrors
}

func getErrorMessage(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return "This field is required"
    case "email":
        return "Invalid email format"
    case "min":
        return fmt.Sprintf("Minimum length is %s", err.Param())
    case "max":
        return fmt.Sprintf("Maximum length is %s", err.Param())
    case "gte":
        return fmt.Sprintf("Must be >= %s", err.Param())
    case "lte":
        return fmt.Sprintf("Must be <= %s", err.Param())
    default:
        return "Invalid value"
    }
}

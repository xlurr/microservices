package models

import "time"

type User struct {
    ID        int64     `json:"id"`
    FirstName string    `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string    `json:"last_name" validate:"required,min=2,max=50"`
    Email     string    `json:"email" validate:"required,email"`
    Age       int       `json:"age" validate:"required,gte=18,lte=120"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
    FirstName string `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string `json:"last_name" validate:"required,min=2,max=50"`
    Email     string `json:"email" validate:"required,email"`
    Age       int    `json:"age" validate:"required,gte=18,lte=120"`
}

type ErrorResponse struct {
    Status  int          `json:"status"`
    Message string       `json:"message"`
    Errors  []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

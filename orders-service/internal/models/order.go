package models

import "time"

type Order struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id" validate:"required"`
    Items       []string  `json:"items" validate:"required,min=1"`
    TotalAmount float64   `json:"total_amount"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
    UserID int64    `json:"user_id" validate:"required"`
    Items  []string `json:"items" validate:"required,min=1"`
}

type UpdateStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=created payment_pending paid delivered cancelled"`
}

type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}

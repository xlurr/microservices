package models

import "time"

type Payment struct {
    ID        int64     `json:"id"`
    OrderID   int64     `json:"order_id"`
    Amount    float64   `json:"amount"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

type CreatePaymentRequest struct {
    OrderID int64   `json:"order_id" validate:"required"`
    Amount  float64 `json:"amount" validate:"required,gt=0"`
}

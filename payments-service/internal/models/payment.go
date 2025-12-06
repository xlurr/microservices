package models

import "time"

type Payment struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	OrderID   int64     `json:"orderId"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // pending, completed, failed
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

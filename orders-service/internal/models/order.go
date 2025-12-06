package models

import "time"

type Order struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Items       []string  `json:"items"`
	TotalAmount float64   `json:"totalAmount"`
	Status      string    `json:"status"` // created, processing, completed
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

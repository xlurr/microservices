package models

import "time"

type Delivery struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userId"`
	OrderID    int64     `json:"orderId"`
	Address    string    `json:"address"`
	Status     string    `json:"status"` // pending, shipped, delivered
	TrackingID string    `json:"trackingId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

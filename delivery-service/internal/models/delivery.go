package models

import "time"

type Delivery struct {
    ID        int64     `json:"id"`
    OrderID   int64     `json:"order_id"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

package models

import "time"

type Order struct {
	ID              int64
	UserID          int64
	OrderDetailID   int64
	Amount          float64
	TotalQty        int
	Status          int
	PaymentMethod   string
	ShippingAddress string
}

type OrderDetail struct {
	ID           int64
	Products     string
	OrderHistory string
}

type CheckoutItem struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CheckoutRequest struct {
	UserID           int64          `json:"user_id"`
	Items            []CheckoutItem `json:"items"`
	PaymentMethod    string         `json:"payment_method"`
	ShippingAddress  string         `json:"shipping_address"`
	IdempotencyToken string         `json:"idempotency_token"`
}
type OrderHistoryParam struct {
	UserID int64
	Status int
}
type OrderHistoryResponse struct {
	OrderID         int64           `jsonn:"order_id"`
	TotalAmount     float64         `json:"total_amount"`
	TotalQty        int             `json:"total_qty"`
	Status          string          `json:"status"`
	PaymentMethod   string          `json:"payment_method"`
	ShippingAddress string          `json:"shipping_address"`
	Products        []CheckoutItem  `json:"products"`
	History         []StatusHistory `json:"history"`
}

type StatusHistory struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type OrderRequestLog struct {
	ID               int64     `json:"id"`
	IdempotencyToken string    `json:"idempotency_token"`
	CreateTime       time.Time `json:"create_time"`
}

type OrderHistoryResult struct {
	ID              int64 `gorm:"column.id"`
	Amount          float64
	TotalQuantity   int `gorm:"column:total_qty"`
	Status          int
	PaymentMethod   string
	ShippingAddress string
	Products        string `gorm:"column:products"`
	OrderHistory    string `gorm:"column:order_history"`
}

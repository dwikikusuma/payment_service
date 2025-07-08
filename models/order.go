package models

type OrderCreatedEvent struct {
	OrderID         int64   `json:"order_id"`
	UserID          int64   `json:"user_id"`
	Amount          float64 `json:"amount"`
	PaymentMethod   string  `json:"payment_method"`
	PayerEmail      string  `json:"payer_email"`
	ShippingAddress string  `json:"shipping_address"`
}

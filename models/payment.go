package models

import "time"

type Payment struct {
	ID          int64     `json:"id"`
	OrderID     int64     `json:"order_id"`
	UserID      int64     `json:"user_id"`
	ExternalID  string    `json:"external_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	CreateTime  time.Time `json:"create_time"`
	ExpiredTime time.Time `json:"expired_time"`
	UpdateTime  time.Time `json:"update_time"`
}

type Status struct {
	ID   int64
	Name string
}

type PaymentRequests struct {
	ID         int64     `json:"id"`
	OrderID    int64     `json:"order_id"`
	UserID     int64     `json:"user_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	RetryCount int       `json:"retry_count"`
	Notes      string    `json:"notes"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// TableName overrides the default table name used by GORM.
//
// By default, GORM uses the pluralized form of the struct name as the table name.
// Defining this method allows you to explicitly set a custom table name.
//
// For example, the Payment struct will map to the "payments" table in the database.
func (Payment) TableName() string {
	return "payments"
}

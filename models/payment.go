package models

type Payment struct {
	ID         int64   `json:"id"`
	OrderID    int64   `json:"order_id"`
	UserID     int64   `json:"user_id"`
	ExternalID string  `json:"external_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

type Status struct {
	ID   int64
	Name string
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

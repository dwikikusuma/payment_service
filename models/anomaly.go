package models

import "time"

type PaymentAnomaly struct {
	ID          int       `json:"id"`
	OrderID     int64     `json:"order_id"`
	ExternalID  string    `json:"external_id"`
	AnomalyType int       `json:"anomaly_type"` // 1: Anomaly Amount
	Notes       string    `json:"notes"`
	Status      int       `json:"status"` // 1: Success, 99: Need to check
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func (PaymentAnomaly) TableName() string {
	return "payment_anomalies"
}

package models

import "time"

type FailedEvents struct {
	ID         int       `json:"id"`
	OrderID    int64     `json:"order_id"`
	ExternalID string    `json:"external_id"`
	FailedType int       `json:"failed_type"`
	Notes      string    `json:"notes"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (FailedEvents) TableName() string {
	return "failed_events"
}

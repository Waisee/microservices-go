package record

import "time"

type OrderRecord struct {
	UUID            string
	Items           []OrderItemRecord
	TransactionUUID string
	PaymentMethod   string
	Status          string
	CreatedAt       time.Time
}

type OrderItemRecord struct {
	PartUUID string
	PartType string
	Price    int64
}

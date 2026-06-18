package record

import "time"

type OrderRecord struct {
	UUID            string            `db:"uuid"`
	Items           []OrderItemRecord `db:"items"`
	TransactionUUID string            `db:"transaction_uuid"`
	PaymentMethod   string            `db:"payment_method"`
	Status          string            `db:"status"`
	CreatedAt       time.Time         `db:"created_at"`
}

type OrderItemRecord struct {
	PartUUID string `db:"part_uuid"`
	PartType string `db:"part_type"`
	Price    int64  `db:"price"`
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UUID            uuid.UUID
	Items           []OrderItem
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// TotalPrice возвращает сумму цен всех позиций заказа.
func (o Order) TotalPrice() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.Price
	}
	return total
}

type OrderItem struct {
	PartUUID uuid.UUID
	PartType PartType
	Price    int64
}

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

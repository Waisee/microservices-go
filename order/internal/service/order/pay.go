package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
)

func (s *OrderService) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.OrderRepository.Get(ctx, orderUUID)
	if err != nil {
		return uuid.UUID{}, err
	}
	switch order.Status {
	case model.OrderStatusPendingPayment:
	case model.OrderStatusPaid:
		return uuid.UUID{}, errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return uuid.UUID{}, errs.ErrOrderCancelled
	default:
		return uuid.UUID{}, errs.ErrOrderAlreadyPaid
	}

	transactionUUID, err := s.PaymentClient.PayOrder(ctx, orderUUID, method)
	if err != nil {
		return uuid.UUID{}, err
	}
	order.TransactionUUID = &transactionUUID
	paymentMethod := method
	order.PaymentMethod = &paymentMethod
	order.Status = model.OrderStatusPaid

	if err := s.OrderRepository.Update(ctx, order); err != nil {
		return uuid.UUID{}, err
	}
	return transactionUUID, nil
}

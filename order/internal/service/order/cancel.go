package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
)

func (s *OrderService) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.OrderRepository.Get(ctx, orderUUID)
	if err != nil {
		return errs.ErrOrderNotFound
	}
	switch order.Status {
	case model.OrderStatusPendingPayment:
	case model.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return errs.ErrOrderCancelled
	default:
		return errs.ErrOrderCancelled
	}
	order.Status = model.OrderStatusCancelled
	if err := s.OrderRepository.Update(ctx, order); err != nil {
		return err
	}
	return nil
}

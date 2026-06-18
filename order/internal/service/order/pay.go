package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
)

func (s *OrderService) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	var transactionUUID uuid.UUID

	err := s.TxManager.Do(ctx, func(ctx context.Context) error {
		// 1. Получение заказа.
		order, err := s.OrderRepository.Get(ctx, orderUUID)
		if err != nil {
			return err
		}
		// 2. Проверка статуса заказа.
		switch order.Status {
		case model.OrderStatusPendingPayment:
		case model.OrderStatusPaid:
			return errs.ErrOrderAlreadyPaid
		case model.OrderStatusCancelled:
			return errs.ErrOrderCancelled
		default:
			return errs.ErrOrderAlreadyPaid
		}

		// 3. Оплата заказа.
		transactionUUID, err = s.PaymentClient.PayOrder(ctx, orderUUID, method)
		if err != nil {
			return err
		}

		order.TransactionUUID = &transactionUUID
		paymentMethod := method
		order.PaymentMethod = &paymentMethod
		order.Status = model.OrderStatusPaid

		// 4. Обновление заказа.
		if err := s.OrderRepository.Update(ctx, order); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return transactionUUID, nil
}

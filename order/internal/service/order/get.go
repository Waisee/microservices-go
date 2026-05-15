package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/model"
)

func (s *OrderService) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.OrderRepository.Get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

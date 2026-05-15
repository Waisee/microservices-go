package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/order/internal/service/input"
)

func (s *OrderService) Create(ctx context.Context, in input.CreateOrderInput) (model.Order, error) {
	items := make([]model.OrderItem, 0, 4)
	parts, err := s.InventoryClient.ListParts(ctx, in.PartUUIDs())
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return model.Order{}, errs.ErrPartNotFound
		}
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			return model.Order{}, errs.ErrInvalidPartUUID
		}
		return model.Order{}, err
	}
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, errs.ErrOutOfStock
		}
		items = append(items, model.OrderItem{
			PartUUID: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		})
	}
	order := model.Order{
		UUID:      uuid.New(),
		Items:     items,
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now(),
	}

	if err := s.OrderRepository.Create(ctx, order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

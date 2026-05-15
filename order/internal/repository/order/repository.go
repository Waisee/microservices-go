package order

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/order/internal/repository/converter"
	"github.com/waisee/microservices-go/order/internal/repository/record"
)

type OrderRepository struct {
	orders map[string]record.OrderRecord
	mu     sync.RWMutex
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]record.OrderRecord),
		mu:     sync.RWMutex{},
	}
}

func (r *OrderRepository) Create(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.UUID.String()] = converter.ModelToRecord(order)

	return nil
}

func (r *OrderRepository) Get(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[uuid.String()]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}
	return converter.RecordToModel(order), nil
}

func (r *OrderRepository) Update(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.UUID.String()] = converter.ModelToRecord(order)

	return nil
}

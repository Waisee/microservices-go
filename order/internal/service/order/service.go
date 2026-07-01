package order

import "context"

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type OrderService struct {
	OrderRepository OrderRepository
	PaymentClient   PaymentClient
	InventoryClient InventoryClient
	TxManager       TxManager
}

func NewOrderService(
	orderRepository OrderRepository,
	paymentClient PaymentClient,
	inventoryClient InventoryClient,
	txManager TxManager,
) *OrderService {
	return &OrderService{
		OrderRepository: orderRepository,
		PaymentClient:   paymentClient,
		InventoryClient: inventoryClient,
		TxManager:       txManager,
	}
}

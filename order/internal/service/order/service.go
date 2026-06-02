package order

import "sync"

type OrderService struct {
	OrderRepository OrderRepository
	PaymentClient   PaymentClient
	InventoryClient InventoryClient

	mu sync.Mutex
}

func NewOrderService(orderRepository OrderRepository, paymentClient PaymentClient, inventoryClient InventoryClient) *OrderService {
	return &OrderService{OrderRepository: orderRepository, PaymentClient: paymentClient, InventoryClient: inventoryClient}
}

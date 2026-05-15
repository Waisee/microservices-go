package order

type OrderService struct {
	OrderRepository OrderRepository
	PaymentClient   PaymentClient
	InventoryClient InventoryClient
}

func NewOrderService(orderRepository OrderRepository, paymentClient PaymentClient, inventoryClient InventoryClient) *OrderService {
	return &OrderService{OrderRepository: orderRepository, PaymentClient: paymentClient, InventoryClient: inventoryClient}
}

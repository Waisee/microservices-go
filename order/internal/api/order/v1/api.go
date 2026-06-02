package v1

import orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"

type OrderAPI struct {
	orderv1.UnimplementedHandler
	orderService OrderService
}

func NewOrderAPI(orderService OrderService) *OrderAPI {
	return &OrderAPI{orderService: orderService}
}

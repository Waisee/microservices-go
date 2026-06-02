package app

import (
	"log/slog"
	"net/http"

	orderapi "github.com/waisee/microservices-go/order/internal/api/order/v1"
	inventoryclientv1 "github.com/waisee/microservices-go/order/internal/clients/grpc/inventory/v1"
	paymentclientv1 "github.com/waisee/microservices-go/order/internal/clients/grpc/payment/v1"
	orderrepo "github.com/waisee/microservices-go/order/internal/repository/order"
	ordersvc "github.com/waisee/microservices-go/order/internal/service/order"
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"
)

func NewHTTPHandler(inventoryServiceClient inventoryv1.InventoryServiceClient, paymentServiceClient paymentv1.PaymentServiceClient) (http.Handler, error) {
	inventoryClient := inventoryclientv1.NewInventoryClient(
		inventoryServiceClient,
	)

	paymentClient := paymentclientv1.NewPaymentClient(
		paymentServiceClient,
	)

	repo := orderrepo.NewOrderRepository()
	service := ordersvc.NewOrderService(repo, paymentClient, inventoryClient)
	api := orderapi.NewOrderAPI(service)

	return orderv1.NewServer(api, orderv1.WithErrorHandler(orderapi.OgenErrorHandler(slog.Default())))
}

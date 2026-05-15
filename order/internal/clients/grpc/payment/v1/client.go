package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/clients/grpc/payment/v1/converter"
	"github.com/waisee/microservices-go/order/internal/model"
	paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"
)

type PaymentClient struct {
	paymentv1.PaymentServiceClient
}

func NewPaymentClient(grpc paymentv1.PaymentServiceClient) *PaymentClient {
	return &PaymentClient{
		PaymentServiceClient: grpc,
	}
}

func (c *PaymentClient) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	resp, err := c.PaymentServiceClient.PayOrder(ctx, converter.ToPayOrderRequest(orderUUID, method))
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.MustParse(resp.GetTransactionUuid()), nil
}

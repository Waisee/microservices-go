package v1

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/clients/grpc/payment/v1/converter"
	"github.com/waisee/microservices-go/order/internal/model"
	paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"
)

const (
	timeout = 5 * time.Second
)

type PaymentClient struct {
	client paymentv1.PaymentServiceClient
}

func NewPaymentClient(c paymentv1.PaymentServiceClient) *PaymentClient {
	return &PaymentClient{
		client: c,
	}
}

func (c *PaymentClient) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := c.client.PayOrder(ctx, converter.ToPayOrderRequest(orderUUID, method))
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.MustParse(resp.GetTransactionUuid()), nil
}

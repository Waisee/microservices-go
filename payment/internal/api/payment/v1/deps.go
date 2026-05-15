package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/payment/internal/service/input"
)

type PaymentService interface {
	Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error)
}

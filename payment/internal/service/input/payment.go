package input

import (
	"github.com/google/uuid"

	"github.com/waisee/microservices-go/payment/internal/model"
)

type PayOrderInput struct {
	OrderUUID     uuid.UUID
	PaymentMethod model.PaymentMethod
}

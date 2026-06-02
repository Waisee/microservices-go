package converter

import (
	"github.com/google/uuid"

	"github.com/waisee/microservices-go/payment/internal/errors"
	"github.com/waisee/microservices-go/payment/internal/model"
	"github.com/waisee/microservices-go/payment/internal/service/input"
	paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"
)

func ProtoToPayOrderInput(in *paymentv1.PayOrderRequest) (input.PayOrderInput, error) {
	orderUUID, err := uuid.Parse(in.GetOrderUuid())
	if err != nil || orderUUID == uuid.Nil {
		return input.PayOrderInput{}, errs.ErrInvalidOrderUUID
	}
	method, err := paymentMethodToModel(in.GetPaymentMethod())
	if err != nil {
		return input.PayOrderInput{}, err
	}
	return input.PayOrderInput{
		OrderUUID:     orderUUID,
		PaymentMethod: method,
	}, nil
}

func paymentMethodToModel(m paymentv1.PaymentMethod) (model.PaymentMethod, error) {
	switch m {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		return "", errs.ErrInvalidPaymentMethod
	default:
		return "", errs.ErrInvalidPaymentMethod
	}
}

func UUIDToProtoPayOrderResponse(in uuid.UUID) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: in.String(),
	}
}

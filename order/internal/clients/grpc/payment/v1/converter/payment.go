package converter

import (
	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/model"
	paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"
)

func ToPayOrderRequest(orderUUID uuid.UUID, method model.PaymentMethod) *paymentv1.PayOrderRequest {
	return &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: paymentMethodToProto(method),
	}
}

func paymentMethodToProto(method model.PaymentMethod) paymentv1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	}
	return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
}

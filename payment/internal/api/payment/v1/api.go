package v1

import paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"

type PaymentAPI struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

func NewPaymentAPI(paymentService PaymentService) *PaymentAPI {
	return &PaymentAPI{
		paymentService: paymentService,
	}
}

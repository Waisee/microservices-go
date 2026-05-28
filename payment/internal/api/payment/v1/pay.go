package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/waisee/microservices-go/payment/internal/api/converter"
	"github.com/waisee/microservices-go/payment/internal/errors"
	paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"
)

// PayOrder проводит оплату заказа и возвращает UUID транзакции.
// Возвращает InvalidArgument при невалидном UUID заказа или способе оплаты.
func (a *PaymentAPI) PayOrder(ctx context.Context, in *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	input, err := converter.ProtoToPayOrderInput(in)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidOrderUUID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}

	transactionUUID, err := a.paymentService.Pay(ctx, input)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}

	return converter.UUIDToProtoPayOrderResponse(transactionUUID), nil
}

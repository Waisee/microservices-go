package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/waisee/microservices-go/order/internal/api/order/v1/converter"
	"github.com/waisee/microservices-go/order/internal/errors"
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

// PayOrder оплачивает заказ выбранным способом оплаты.
// Возвращает 404, если заказ не найден, 409 если заказ уже оплачен или отменён.
func (api *OrderAPI) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	transactionUUID, err := api.orderService.Pay(ctx, params.OrderUUID, converter.ProtoToPaymentMethod(req))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.PayOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case errors.Is(err, errs.ErrOrderAlreadyPaid),
			errors.Is(err, errs.ErrOrderCancelled):
			return &orderv1.PayOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		default:
			return &orderv1.PayOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка",
			}, nil
		}
	}
	proto, err := converter.ModelToPayOrderRes(transactionUUID)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

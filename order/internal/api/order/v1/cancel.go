package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/waisee/microservices-go/order/internal/errors"
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

func (api *OrderAPI) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := api.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.CancelOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case errors.Is(err, errs.ErrOrderAlreadyPaid),
			errors.Is(err, errs.ErrOrderCancelled):
			return &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		default:
			return &orderv1.CancelOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка",
			}, nil
		}
	}
	return &orderv1.CancelOrderResponse{}, nil
}

package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/waisee/microservices-go/order/internal/api/order/v1/converter"
	"github.com/waisee/microservices-go/order/internal/errors"
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

// GetOrder возвращает заказ по UUID из пути запроса.
// Возвращает 404, если заказ не найден.
func (api *OrderAPI) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := api.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.GetOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		default:
			return &orderv1.GetOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка",
			}, nil
		}
	}
	proto, err := converter.ModelToGetOrderRes(order)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

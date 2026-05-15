package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/waisee/microservices-go/order/internal/api/order/v1/converter"
	"github.com/waisee/microservices-go/order/internal/errors"
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

func (api *OrderAPI) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	order, err := api.orderService.Create(ctx, converter.ProtoToCreateOrderInput(req))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrPartNotFound):
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case errors.Is(err, errs.ErrOutOfStock):
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		default:
			return &orderv1.CreateOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка",
			}, nil
		}
	}
	proto, err := converter.ModelToCreateOrderRes(order)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

package v1_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderv1api "github.com/waisee/microservices-go/order/internal/api/order/v1"
	"github.com/waisee/microservices-go/order/internal/api/order/v1/mocks"
	errs "github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/order/internal/service/input"
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx        = context.Background()
		hullUUID   = uuid.New()
		engineUUID = uuid.New()
		orderUUID  = uuid.New()
		repoErr    = gofakeit.Error()

		order = model.Order{
			UUID: orderUUID,
			Items: []model.OrderItem{
				{PartUUID: hullUUID, PartType: model.PartTypeHull, Price: 500000},
				{PartUUID: engineUUID, PartType: model.PartTypeEngine, Price: 300000},
			},
			Status: model.OrderStatusPendingPayment,
		}

		req = &orderv1.CreateOrderRequest{
			HullUUID:   hullUUID,
			EngineUUID: engineUUID,
		}
	)

	tests := []struct {
		name      string
		setupMock func(svc *mocks.OrderService)
		assertRes func(t *testing.T, res orderv1.CreateOrderRes, err error)
	}{
		{
			name: "успешное создание заказа",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderInput{HullUUID: hullUUID, EngineUUID: engineUUID}).
					Return(order, nil)
			},
			assertRes: func(t *testing.T, res orderv1.CreateOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.CreateOrderResponse)
				require.True(t, ok)
				assert.Equal(t, orderUUID, resp.OrderUUID)
				assert.Equal(t, int64(800000), resp.TotalPrice)
			},
		},
		{
			name: "деталь не найдена",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderInput{HullUUID: hullUUID, EngineUUID: engineUUID}).
					Return(model.Order{}, errs.ErrPartNotFound)
			},
			assertRes: func(t *testing.T, res orderv1.CreateOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.CreateOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "деталь отсутствует на складе",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderInput{HullUUID: hullUUID, EngineUUID: engineUUID}).
					Return(model.Order{}, errs.ErrOutOfStock)
			},
			assertRes: func(t *testing.T, res orderv1.CreateOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.CreateOrderConflict)
				require.True(t, ok)
				assert.Equal(t, http.StatusConflict, resp.Code)
			},
		},
		{
			name: "внутренняя ошибка",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Create(ctx, input.CreateOrderInput{HullUUID: hullUUID, EngineUUID: engineUUID}).
					Return(model.Order{}, repoErr)
			},
			assertRes: func(t *testing.T, res orderv1.CreateOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.CreateOrderInternalServerError)
				require.True(t, ok)
				assert.Equal(t, http.StatusInternalServerError, resp.Code)
				assert.Equal(t, "внутренняя ошибка", resp.Message)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewOrderService(t)
			tc.setupMock(svc)

			api := orderv1api.NewOrderAPI(svc)
			res, err := api.CreateOrder(ctx, req)
			tc.assertRes(t, res, err)
		})
	}
}

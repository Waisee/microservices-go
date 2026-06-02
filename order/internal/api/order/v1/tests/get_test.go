package v1_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	orderv1api "github.com/waisee/microservices-go/order/internal/api/order/v1"
	"github.com/waisee/microservices-go/order/internal/api/order/v1/mocks"
	errs "github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

func TestGetOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		repoErr    = gofakeit.Error()
		orderUUID  = uuid.New()
		hullUUID   = uuid.New()
		engineUUID = uuid.New()

		order = model.Order{
			UUID: orderUUID,
			Items: []model.OrderItem{
				{PartUUID: hullUUID, PartType: model.PartTypeHull, Price: 500000},
				{PartUUID: engineUUID, PartType: model.PartTypeEngine, Price: 300000},
			},
			Status:    model.OrderStatusPendingPayment,
			CreatedAt: time.Now(),
		}

		params = orderv1.GetOrderParams{OrderUUID: orderUUID}
	)

	tests := []struct {
		name      string
		setupMock func(svc *mocks.OrderService)
		assertRes func(t *testing.T, res orderv1.GetOrderRes, err error)
	}{
		{
			name: "заказ найден",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().Get(ctx, orderUUID).Return(order, nil)
			},
			assertRes: func(t *testing.T, res orderv1.GetOrderRes, err error) {
				require.NoError(t, err)
				dto, ok := res.(*orderv1.OrderDto)
				require.True(t, ok)
				assert.Equal(t, orderUUID, dto.OrderUUID)
				assert.Equal(t, hullUUID, dto.HullUUID)
				assert.Equal(t, engineUUID, dto.EngineUUID)
				assert.Equal(t, int64(800000), dto.TotalPrice)
				assert.Equal(t, orderv1.OrderStatusPENDINGPAYMENT, dto.Status)
			},
		},
		{
			name: "заказ не найден",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().Get(ctx, orderUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			assertRes: func(t *testing.T, res orderv1.GetOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.GetOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "внутренняя ошибка",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().Get(ctx, orderUUID).Return(model.Order{}, repoErr)
			},
			assertRes: func(t *testing.T, res orderv1.GetOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.GetOrderInternalServerError)
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
			res, err := api.GetOrder(ctx, params)
			tc.assertRes(t, res, err)
		})
	}
}

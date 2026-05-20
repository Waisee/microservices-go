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
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		repoErr   = gofakeit.Error()

		params = orderv1.CancelOrderParams{OrderUUID: orderUUID}
	)

	tests := []struct {
		name      string
		setupMock func(svc *mocks.OrderService)
		assertRes func(t *testing.T, res orderv1.CancelOrderRes, err error)
	}{
		{
			name: "успешная отмена",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().Cancel(ctx, orderUUID).Return(nil)
			},
			assertRes: func(t *testing.T, res orderv1.CancelOrderRes, err error) {
				require.NoError(t, err)
				_, ok := res.(*orderv1.CancelOrderResponse)
				require.True(t, ok)
			},
		},
		{
			name: "заказ не найден",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().Cancel(ctx, orderUUID).Return(errs.ErrOrderNotFound)
			},
			assertRes: func(t *testing.T, res orderv1.CancelOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.CancelOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "заказ уже оплачен",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().Cancel(ctx, orderUUID).Return(errs.ErrOrderAlreadyPaid)
			},
			assertRes: func(t *testing.T, res orderv1.CancelOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.CancelOrderConflict)
				require.True(t, ok)
				assert.Equal(t, http.StatusConflict, resp.Code)
			},
		},
		{
			name: "заказ отменён",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().Cancel(ctx, orderUUID).Return(errs.ErrOrderCancelled)
			},
			assertRes: func(t *testing.T, res orderv1.CancelOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.CancelOrderConflict)
				require.True(t, ok)
				assert.Equal(t, http.StatusConflict, resp.Code)
			},
		},
		{
			name: "внутренняя ошибка",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().Cancel(ctx, orderUUID).Return(repoErr)
			},
			assertRes: func(t *testing.T, res orderv1.CancelOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.CancelOrderInternalServerError)
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
			res, err := api.CancelOrder(ctx, params)
			tc.assertRes(t, res, err)
		})
	}
}

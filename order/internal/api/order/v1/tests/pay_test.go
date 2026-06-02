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
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx = context.Background()

		paymentErr        = gofakeit.Error()
		orderUUID         = uuid.New()
		transactionUUID   = uuid.New()
		paymentMethodCard = model.PaymentMethodCard

		req = &orderv1.PayOrderRequest{
			PaymentMethod: orderv1.PaymentMethodCARD,
		}
		params = orderv1.PayOrderParams{OrderUUID: orderUUID}
	)

	tests := []struct {
		name      string
		setupMock func(svc *mocks.OrderService)
		assertRes func(t *testing.T, res orderv1.PayOrderRes, err error)
	}{
		{
			name: "успешная оплата",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, paymentMethodCard).
					Return(transactionUUID, nil)
			},
			assertRes: func(t *testing.T, res orderv1.PayOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.PayOrderResponse)
				require.True(t, ok)
				assert.Equal(t, transactionUUID, resp.TransactionUUID)
			},
		},
		{
			name: "заказ не найден",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, paymentMethodCard).
					Return(uuid.Nil, errs.ErrOrderNotFound)
			},
			assertRes: func(t *testing.T, res orderv1.PayOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.PayOrderNotFound)
				require.True(t, ok)
				assert.Equal(t, http.StatusNotFound, resp.Code)
			},
		},
		{
			name: "заказ уже оплачен",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, paymentMethodCard).
					Return(uuid.Nil, errs.ErrOrderAlreadyPaid)
			},
			assertRes: func(t *testing.T, res orderv1.PayOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.PayOrderConflict)
				require.True(t, ok)
				assert.Equal(t, http.StatusConflict, resp.Code)
			},
		},
		{
			name: "заказ отменён",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, paymentMethodCard).
					Return(uuid.Nil, errs.ErrOrderCancelled)
			},
			assertRes: func(t *testing.T, res orderv1.PayOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.PayOrderConflict)
				require.True(t, ok)
				assert.Equal(t, http.StatusConflict, resp.Code)
			},
		},
		{
			name: "внутренняя ошибка",
			setupMock: func(svc *mocks.OrderService) {
				svc.EXPECT().
					Pay(ctx, orderUUID, paymentMethodCard).
					Return(uuid.Nil, paymentErr)
			},
			assertRes: func(t *testing.T, res orderv1.PayOrderRes, err error) {
				require.NoError(t, err)
				resp, ok := res.(*orderv1.PayOrderInternalServerError)
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
			res, err := api.PayOrder(ctx, req, params)
			tc.assertRes(t, res, err)
		})
	}
}

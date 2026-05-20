package v1_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentv1api "github.com/waisee/microservices-go/payment/internal/api/payment/v1"
	"github.com/waisee/microservices-go/payment/internal/api/payment/v1/mocks"
	errs "github.com/waisee/microservices-go/payment/internal/errors"
	"github.com/waisee/microservices-go/payment/internal/model"
	"github.com/waisee/microservices-go/payment/internal/service/input"
	paymentv1 "github.com/waisee/microservices-go/shared/pkg/proto/payment/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		orderUUID       = uuid.New()
		transactionUUID = uuid.New()
		serviceErr      = gofakeit.Error()

		validReq = &paymentv1.PayOrderRequest{
			OrderUuid:     orderUUID.String(),
			PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
		}
	)

	tests := []struct {
		name      string
		req       *paymentv1.PayOrderRequest
		setupMock func(svc *mocks.PaymentService)
		assertRes func(t *testing.T, res *paymentv1.PayOrderResponse, err error)
	}{
		{
			name: "успешная оплата",
			req:  validReq,
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					Pay(ctx, input.PayOrderInput{
						OrderUUID:     orderUUID,
						PaymentMethod: model.PaymentMethodCard,
					}).
					Return(transactionUUID, nil)
			},
			assertRes: func(t *testing.T, res *paymentv1.PayOrderResponse, err error) {
				require.NoError(t, err)
				assert.Equal(t, transactionUUID.String(), res.TransactionUuid)
			},
		},
		{
			name: "неверный UUID заказа",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     "invalid-uuid",
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
			},
			setupMock: func(svc *mocks.PaymentService) {},
			assertRes: func(t *testing.T, res *paymentv1.PayOrderResponse, err error) {
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "неверный метод оплаты в запросе",
			req: &paymentv1.PayOrderRequest{
				OrderUuid:     orderUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			},
			setupMock: func(svc *mocks.PaymentService) {},
			assertRes: func(t *testing.T, res *paymentv1.PayOrderResponse, err error) {
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "неверный метод оплаты от сервиса",
			req:  validReq,
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					Pay(ctx, input.PayOrderInput{OrderUUID: orderUUID, PaymentMethod: model.PaymentMethodCard}).
					Return(uuid.Nil, errs.ErrInvalidPaymentMethod)
			},
			assertRes: func(t *testing.T, res *paymentv1.PayOrderResponse, err error) {
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "внутренняя ошибка",
			req:  validReq,
			setupMock: func(svc *mocks.PaymentService) {
				svc.EXPECT().
					Pay(ctx, input.PayOrderInput{OrderUUID: orderUUID, PaymentMethod: model.PaymentMethodCard}).
					Return(uuid.Nil, serviceErr)
			},
			assertRes: func(t *testing.T, res *paymentv1.PayOrderResponse, err error) {
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.Internal, st.Code())
				assert.Equal(t, "внутренняя ошибка", st.Message())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewPaymentService(t)
			tc.setupMock(svc)

			api := paymentv1api.NewPaymentAPI(svc)
			res, err := api.PayOrder(ctx, tc.req)
			switch tc.name {
			case "неверный UUID заказа", "неверный метод оплаты в запросе":
				svc.AssertNotCalled(t, "Pay")
			}
			tc.assertRes(t, res, err)
		})
	}
}

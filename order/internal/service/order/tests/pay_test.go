package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
	orderservice "github.com/waisee/microservices-go/order/internal/service/order"
	"github.com/waisee/microservices-go/order/internal/service/order/mocks"
)

func TestPay(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID     uuid.UUID
		paymentMethod model.PaymentMethod
	}

	type expected struct {
		transactionUUID uuid.UUID
		err             error
	}

	var (
		ctx = context.Background()

		paymentErr = gofakeit.Error()
		repoErr    = gofakeit.Error()

		orderUUID       = uuid.New()
		transactionUUID = uuid.New()
		hullUUID        = uuid.New()
		engineUUID      = uuid.New()

		order = model.Order{
			UUID: orderUUID,
			Items: []model.OrderItem{
				{PartUUID: hullUUID, PartType: model.PartTypeHull, Price: 500000},
				{PartUUID: engineUUID, PartType: model.PartTypeEngine, Price: 300000},
			},
			TransactionUUID: nil,
			PaymentMethod:   nil,
			Status:          model.OrderStatusPendingPayment,
			CreatedAt:       time.Now(),
		}
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, client *mocks.PaymentClient)
		expected  expected
	}{
		{
			name: "успешная оплата",
			args: args{
				orderUUID:     orderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(order, nil)

				client.EXPECT().PayOrder(ctx, orderUUID, model.PaymentMethodCard).Return(transactionUUID, nil)

				repo.EXPECT().Update(ctx, mock.MatchedBy(func(o model.Order) bool {
					return o.Status == model.OrderStatusPaid &&
						o.TransactionUUID != nil &&
						*o.TransactionUUID == transactionUUID &&
						o.PaymentMethod != nil &&
						*o.PaymentMethod == model.PaymentMethodCard
				})).Return(nil)
			},
			expected: expected{
				transactionUUID: transactionUUID,
				err:             nil,
			},
		},
		{
			name: "заказ не найден",
			args: args{
				orderUUID:     orderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				err:             errs.ErrOrderNotFound,
			},
		},
		{
			name: "заказ уже оплачен",
			args: args{
				orderUUID:     orderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.PaymentClient) {
				paidOrder := order
				paidOrder.Status = model.OrderStatusPaid
				repo.EXPECT().Get(ctx, orderUUID).Return(paidOrder, nil)
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				err:             errs.ErrOrderAlreadyPaid,
			},
		},
		{
			name: "заказ отменён",
			args: args{
				orderUUID:     orderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.PaymentClient) {
				cancelledOrder := order
				cancelledOrder.Status = model.OrderStatusCancelled
				repo.EXPECT().Get(ctx, orderUUID).Return(cancelledOrder, nil)
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				err:             errs.ErrOrderCancelled,
			},
		},
		{
			name: "ошибка Payment service",
			args: args{
				orderUUID:     orderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(order, nil)
				client.EXPECT().PayOrder(ctx, orderUUID, model.PaymentMethodCard).Return(uuid.Nil, paymentErr)
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				err:             paymentErr,
			},
		},
		{
			name: "ошибка при обновении",
			args: args{
				orderUUID:     orderUUID,
				paymentMethod: model.PaymentMethodCard,
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(order, nil)
				client.EXPECT().PayOrder(ctx, orderUUID, model.PaymentMethodCard).Return(transactionUUID, nil)
				repo.EXPECT().Update(ctx, mock.MatchedBy(func(o model.Order) bool {
					return o.Status == model.OrderStatusPaid
				})).Return(repoErr)
			},
			expected: expected{
				transactionUUID: uuid.Nil,
				err:             repoErr,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			paymentClient := mocks.NewPaymentClient(t)
			tc.setupMock(orderRepo, paymentClient)

			svc := orderservice.NewOrderService(orderRepo, paymentClient, nil)
			got, err := svc.Pay(ctx, tc.args.orderUUID, tc.args.paymentMethod)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, uuid.Nil, got)
				switch tc.name {
				case "заказ не найден", "заказ уже оплачен", "заказ отменён":
					paymentClient.AssertNotCalled(t, "PayOrder")
					orderRepo.AssertNotCalled(t, "Update")
				case "ошибка Payment service":
					orderRepo.AssertNotCalled(t, "Update")
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected.transactionUUID, got)
			assert.NotEqual(t, uuid.Nil, got)
		})
	}
}

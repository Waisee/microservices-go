package order_test

import (
	"context"
	"testing"

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

func TestCancel(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
	}

	type expected struct {
		err error
	}

	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
		repoErr   = gofakeit.Error()
		order     = model.Order{
			UUID:   orderUUID,
			Status: model.OrderStatusPendingPayment,
		}
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "успешная отмена",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(order, nil)
				repo.EXPECT().Update(ctx, mock.MatchedBy(func(o model.Order) bool {
					return o.UUID == orderUUID && o.Status == model.OrderStatusCancelled
				})).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "заказ не найден",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				err: errs.ErrOrderNotFound,
			},
		},
		{
			name: "ошибка репозитория",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{}, repoErr)
			},
			expected: expected{
				err: repoErr,
			},
		},
		{
			name: "заказ уже оплачен",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				paidOrder := order
				paidOrder.Status = model.OrderStatusPaid
				repo.EXPECT().Get(ctx, orderUUID).Return(paidOrder, nil)
			},
			expected: expected{
				err: errs.ErrOrderAlreadyPaid,
			},
		},
		{
			name: "заказ отменён",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				cancelledOrder := order
				cancelledOrder.Status = model.OrderStatusCancelled
				repo.EXPECT().Get(ctx, orderUUID).Return(cancelledOrder, nil)
			},
			expected: expected{
				err: errs.ErrOrderCancelled,
			},
		},
		{
			name: "ошибка при обновении",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(order, nil)
				repo.EXPECT().Update(ctx, mock.MatchedBy(func(o model.Order) bool {
					return o.Status == model.OrderStatusCancelled
				})).Return(repoErr)
			},
			expected: expected{
				err: repoErr,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			tc.setupMock(orderRepo)

			svc := orderservice.NewOrderService(orderRepo, nil, nil)
			err := svc.Cancel(ctx, tc.args.orderUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				switch tc.name {
				case "заказ не найден", "ошибка репозитория", "заказ уже оплачен", "заказ отменён":
					orderRepo.AssertNotCalled(t, "Update")
				}
				return
			}
			require.NoError(t, err)
		})
	}
}

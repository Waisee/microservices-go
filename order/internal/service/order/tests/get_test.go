package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
	orderservice "github.com/waisee/microservices-go/order/internal/service/order"
	"github.com/waisee/microservices-go/order/internal/service/order/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
	}

	type expected struct {
		order model.Order
		err   error
	}

	var (
		ctx = context.Background()

		repoErr = gofakeit.Error()

		orderUUID  = uuid.New()
		hullUUID   = uuid.New()
		engineUUID = uuid.New()

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
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "заказ найден",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, orderUUID).Return(order, nil)
			},
			expected: expected{
				order: order,
				err:   nil,
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
				order: model.Order{},
				err:   errs.ErrOrderNotFound,
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
				order: model.Order{},
				err:   repoErr,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)

			tc.setupMock(orderRepo)

			svc := orderservice.NewOrderService(orderRepo, nil, nil)
			got, err := svc.Get(ctx, tc.args.orderUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, got.UUID)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected.order, got)
		})
	}
}

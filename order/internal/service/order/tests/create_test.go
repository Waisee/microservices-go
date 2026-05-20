package order_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/order/internal/service/input"
	orderservice "github.com/waisee/microservices-go/order/internal/service/order"
	"github.com/waisee/microservices-go/order/internal/service/order/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		in input.CreateOrderInput
	}

	type expected struct {
		err error
	}

	var (
		ctx = context.Background()

		repoErr = gofakeit.Error()

		hullUUID   = uuid.New()
		engineUUID = uuid.New()
		shieldUUID = uuid.New()
		weaponUUID = uuid.New()

		partsInStock = []model.Part{
			{UUID: hullUUID, Name: "Hull", PartType: model.PartTypeHull, Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", PartType: model.PartTypeEngine, Price: 300000, StockQuantity: 5},
		}

		partsOutOfStock = []model.Part{
			{UUID: hullUUID, Name: "Hull", PartType: model.PartTypeHull, Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", PartType: model.PartTypeEngine, Price: 300000, StockQuantity: 0},
		}

		partsWithOptional = []model.Part{
			{UUID: hullUUID, Name: "Hull", PartType: model.PartTypeHull, Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", PartType: model.PartTypeEngine, Price: 300000, StockQuantity: 5},
			{UUID: shieldUUID, Name: "Shield", PartType: model.PartTypeShield, Price: 100000, StockQuantity: 3},
			{UUID: weaponUUID, Name: "Weapon", PartType: model.PartTypeWeapon, Price: 200000, StockQuantity: 2},
		}
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, client *mocks.InventoryClient)
		expected  expected
	}{
		{
			name: "успешное создание заказа",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(partsInStock, nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(o model.Order) bool {
						return len(o.Items) == 2 &&
							o.Items[0].PartUUID == hullUUID &&
							o.Items[1].PartUUID == engineUUID &&
							o.TotalPrice() == 800000 && // 500000 + 300000
							o.Status == model.OrderStatusPendingPayment
					})).
					Return(nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "создание с shield и weapon",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
					ShieldUUID: &shieldUUID,
					WeaponUUID: &weaponUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID, shieldUUID, weaponUUID}).
					Return(partsWithOptional, nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(o model.Order) bool {
						return len(o.Items) == 4 &&
							o.Items[0].PartUUID == hullUUID &&
							o.Items[1].PartUUID == engineUUID &&
							o.Items[2].PartUUID == shieldUUID &&
							o.Items[3].PartUUID == weaponUUID &&
							o.TotalPrice() == 1100000 && // 500k+300k+100k+200k
							o.Status == model.OrderStatusPendingPayment
					})).
					Return(nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "деталь не найдена",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
		{
			name: "деталь закончилась на складе",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(partsOutOfStock, nil)
			},
			expected: expected{err: errs.ErrOutOfStock},
		},
		{
			name: "ошибка при сохранении",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []uuid.UUID{hullUUID, engineUUID}).
					Return(partsInStock, nil)
				repo.EXPECT().
					Create(ctx, mock.AnythingOfType("model.Order")).
					Return(repoErr)
			},
			expected: expected{err: repoErr},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			tc.setupMock(orderRepo, inventoryClient)

			svc := orderservice.NewOrderService(orderRepo, paymentClient, inventoryClient)
			order, err := svc.Create(ctx, tc.args.in)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Empty(t, order.UUID)
				switch tc.name {
				case "деталь не найдена", "деталь закончилась на складе":
					orderRepo.AssertNotCalled(t, "Create")
				}
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, order.UUID)
			}
		})
	}
}

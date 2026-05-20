package v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryv1api "github.com/waisee/microservices-go/inventory/internal/api/inventory/v1"
	"github.com/waisee/microservices-go/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/waisee/microservices-go/inventory/internal/errors"
	"github.com/waisee/microservices-go/inventory/internal/model"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
)

func TestGetPart(t *testing.T) {
	t.Parallel()

	var (
		ctx      = context.Background()
		partUUID = gofakeit.UUID()
		repoErr  = gofakeit.Error()

		part = model.Part{
			UUID:          partUUID,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Int64(),
			PartType:      model.PartTypeHull,
			StockQuantity: gofakeit.Int64(),
			CreatedAt:     time.Now(),
		}

		req = &inventoryv1.GetPartRequest{Uuid: partUUID}
	)

	tests := []struct {
		name      string
		setupMock func(svc *mocks.PartService)
		assertRes func(t *testing.T, res *inventoryv1.GetPartResponse, err error)
	}{
		{
			name: "деталь найдена",
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().Get(ctx, partUUID).Return(part, nil)
			},
			assertRes: func(t *testing.T, res *inventoryv1.GetPartResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Part)
				assert.Equal(t, partUUID, res.Part.Uuid)
				assert.Equal(t, part.Name, res.Part.Name)
				assert.Equal(t, inventoryv1.PartType_PART_TYPE_HULL, res.Part.PartType)
			},
		},
		{
			name: "деталь не найдена",
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().Get(ctx, partUUID).Return(model.Part{}, errs.ErrPartNotFound)
			},
			assertRes: func(t *testing.T, res *inventoryv1.GetPartResponse, err error) {
				require.Nil(t, res)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "неверный UUID",
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().Get(ctx, partUUID).Return(model.Part{}, errs.ErrInvalidUUID)
			},
			assertRes: func(t *testing.T, res *inventoryv1.GetPartResponse, err error) {
				require.Nil(t, res)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "внутренняя ошибка",
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().Get(ctx, partUUID).Return(model.Part{}, repoErr)
			},
			assertRes: func(t *testing.T, res *inventoryv1.GetPartResponse, err error) {
				require.Nil(t, res)
				require.Error(t, err)
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

			svc := mocks.NewPartService(t)
			tc.setupMock(svc)

			api := inventoryv1api.NewInventoryApi(svc)
			res, err := api.GetPart(ctx, req)
			tc.assertRes(t, res, err)
		})
	}
}

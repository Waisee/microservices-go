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
	"github.com/waisee/microservices-go/inventory/internal/service/input"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
)

func TestListParts(t *testing.T) {
	t.Parallel()

	var (
		ctx       = context.Background()
		repoErr   = gofakeit.Error()
		missingID = gofakeit.UUID()

		parts = []model.Part{
			{
				UUID:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Int64(),
				PartType:      model.PartTypeHull,
				StockQuantity: gofakeit.Int64(),
				CreatedAt:     time.Now(),
			},
			{
				UUID:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Int64(),
				PartType:      model.PartTypeEngine,
				StockQuantity: gofakeit.Int64(),
				CreatedAt:     time.Now(),
			},
		}
	)

	tests := []struct {
		name      string
		req       *inventoryv1.ListPartsRequest
		setupMock func(svc *mocks.PartService)
		assertRes func(t *testing.T, res *inventoryv1.ListPartsResponse, err error)
	}{
		{
			name: "без фильтра",
			req:  &inventoryv1.ListPartsRequest{},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().List(ctx, input.PartFilter{PartType: input.PartTypeUnspecified}).Return(parts, nil)
			},
			assertRes: func(t *testing.T, res *inventoryv1.ListPartsResponse, err error) {
				require.NoError(t, err)
				require.Len(t, res.Parts, 2)
			},
		},
		{
			name: "фильтр по типу HULL",
			req:  &inventoryv1.ListPartsRequest{PartType: inventoryv1.PartType_PART_TYPE_HULL},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().List(ctx, input.PartFilter{PartType: input.PartTypeHull}).Return(parts[:1], nil)
			},
			assertRes: func(t *testing.T, res *inventoryv1.ListPartsResponse, err error) {
				require.NoError(t, err)
				require.Len(t, res.Parts, 1)
				assert.Equal(t, inventoryv1.PartType_PART_TYPE_HULL, res.Parts[0].PartType)
			},
		},
		{
			name: "деталь не найдена",
			req:  &inventoryv1.ListPartsRequest{Uuids: []string{parts[0].UUID, missingID}},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().List(ctx, input.PartFilter{
					UUIDs:    []string{parts[0].UUID, missingID},
					PartType: input.PartTypeUnspecified,
				}).Return(nil, errs.ErrPartNotFound)
			},
			assertRes: func(t *testing.T, res *inventoryv1.ListPartsResponse, err error) {
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "неверный UUID",
			req:  &inventoryv1.ListPartsRequest{Uuids: []string{"not-a-uuid"}},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().List(ctx, input.PartFilter{
					UUIDs:    []string{"not-a-uuid"},
					PartType: input.PartTypeUnspecified,
				}).Return(nil, errs.ErrInvalidUUID)
			},
			assertRes: func(t *testing.T, res *inventoryv1.ListPartsResponse, err error) {
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "внутренняя ошибка",
			req:  &inventoryv1.ListPartsRequest{},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().List(ctx, input.PartFilter{PartType: input.PartTypeUnspecified}).Return(nil, repoErr)
			},
			assertRes: func(t *testing.T, res *inventoryv1.ListPartsResponse, err error) {
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

			svc := mocks.NewPartService(t)
			tc.setupMock(svc)

			api := inventoryv1api.NewInventoryApi(svc)
			res, err := api.ListParts(ctx, tc.req)
			tc.assertRes(t, res, err)
		})
	}
}

package part_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/waisee/microservices-go/inventory/internal/errors"
	"github.com/waisee/microservices-go/inventory/internal/model"
	"github.com/waisee/microservices-go/inventory/internal/service/input"
	partservice "github.com/waisee/microservices-go/inventory/internal/service/part"
	"github.com/waisee/microservices-go/inventory/internal/service/part/mocks"
)

func TestList(t *testing.T) {
	t.Parallel()

	type args struct {
		filter input.PartFilter
	}
	type expected struct {
		parts []model.Part
		err   error
	}

	repoErr := gofakeit.Error()
	missingUUID := gofakeit.UUID()

	ctx := context.Background()

	parts := []model.Part{
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

	partsHull := []model.Part{
		{
			UUID:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Int64(),
			PartType:      model.PartTypeHull,
			StockQuantity: gofakeit.Int64(),
			CreatedAt:     time.Now(),
		},
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.PartRepository)
		expected  expected
	}{
		{
			name: "без фильтра (все детали)",
			args: args{
				filter: input.PartFilter{},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, input.PartFilter{}).Return(parts, nil)
			},
			expected: expected{
				parts: parts,
				err:   nil,
			},
		},
		{
			name: "фильтр по типу HULL",
			args: args{
				filter: input.PartFilter{
					PartType: input.PartTypeHull,
				},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, input.PartFilter{
					PartType: input.PartTypeHull,
				}).Return(partsHull, nil)
			},
			expected: expected{
				parts: partsHull,
				err:   nil,
			},
		},
		{
			name: "фильтр по конкретным UUID",
			args: args{
				filter: input.PartFilter{
					UUIDs: []string{parts[0].UUID, parts[1].UUID},
				},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, input.PartFilter{
					UUIDs: []string{parts[0].UUID, parts[1].UUID},
				}).Return(parts, nil)
			},
			expected: expected{
				parts: parts,
				err:   nil,
			},
		},
		{
			name: "один из UUID не найден",
			args: args{
				filter: input.PartFilter{
					UUIDs: []string{parts[0].UUID, missingUUID},
				},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, input.PartFilter{
					UUIDs: []string{parts[0].UUID, missingUUID},
				}).Return([]model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{
				parts: []model.Part{},
				err:   errs.ErrPartNotFound,
			},
		},
		{
			name: "невалидный UUID в фильтре",
			args: args{
				filter: input.PartFilter{
					UUIDs: []string{"invalid-uuid"},
				},
			},
			setupMock: func(repo *mocks.PartRepository) {},
			expected: expected{
				parts: nil,
				err:   errs.ErrInvalidUUID,
			},
		},
		{
			name: "ошибка репозитория",
			args: args{
				filter: input.PartFilter{},
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, input.PartFilter{}).Return([]model.Part{}, repoErr)
			},
			expected: expected{
				parts: []model.Part{},
				err:   repoErr,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewPartRepository(t)
			test.setupMock(repo)

			svc := partservice.NewPartService(repo)
			got, err := svc.List(ctx, test.args.filter)

			if test.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, test.expected.err)
				if test.name == "невалидный UUID в фильтре" {
					repo.AssertNotCalled(t, "List")
				}
			} else {
				require.NoError(t, err)

				switch test.name {
				case "без фильтра (все детали)":
					assert.NotEmpty(t, got)
				case "фильтр по типу HULL":
					for _, p := range got {
						assert.Equal(t, model.PartTypeHull, p.PartType)
					}
				case "фильтр по конкретным UUID":
					requested := test.args.filter.UUIDs
					require.Len(t, got, len(requested))
					for i, id := range requested {
						assert.Equal(t, id, got[i].UUID)
					}
				}
			}
		})
	}
}

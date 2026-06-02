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
	partservice "github.com/waisee/microservices-go/inventory/internal/service/part"
	"github.com/waisee/microservices-go/inventory/internal/service/part/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		partUUID string
	}
	type expected struct {
		part model.Part
		err  error
	}

	var (
		ctx      = context.Background()
		partUUID = gofakeit.UUID()
		part     = model.Part{
			UUID:          partUUID,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Int64(),
			PartType:      model.PartTypeHull,
			StockQuantity: gofakeit.Int64(),
			CreatedAt:     time.Now(),
		}

		partNil = model.Part{}

		repoErr = gofakeit.Error()
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.PartRepository)
		expected  expected
	}{
		{
			name: "деталь найдена",
			args: args{
				partUUID: partUUID,
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, partUUID).Return(part, nil)
			},
			expected: expected{
				part: part,
				err:  nil,
			},
		},
		{
			name: "деталь не найдена",
			args: args{
				partUUID: partUUID,
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, partUUID).Return(partNil, errs.ErrPartNotFound)
			},
			expected: expected{
				part: partNil,
				err:  errs.ErrPartNotFound,
			},
		},
		{
			name: "пустой UUID",
			args: args{
				partUUID: "",
			},
			setupMock: func(repo *mocks.PartRepository) {},
			expected: expected{
				part: partNil,
				err:  errs.ErrInvalidUUID,
			},
		},
		{
			name: "невалидный UUID",
			args: args{
				partUUID: "not-a-uuid",
			},
			setupMock: func(repo *mocks.PartRepository) {},
			expected: expected{
				part: partNil,
				err:  errs.ErrInvalidUUID,
			},
		},
		{
			name: "ошибка репозитория",
			args: args{
				partUUID: partUUID,
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, partUUID).Return(partNil, repoErr)
			},
			expected: expected{
				part: partNil,
				err:  repoErr,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewPartRepository(t)
			tc.setupMock(repo)

			svc := partservice.NewPartService(repo)
			got, err := svc.Get(ctx, tc.args.partUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				switch tc.name {
				case "пустой UUID", "невалидный UUID":
					repo.AssertNotCalled(t, "Get")
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.part, got)
			}
		})
	}
}

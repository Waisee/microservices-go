package part

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/waisee/microservices-go/inventory/internal/errors"
	"github.com/waisee/microservices-go/inventory/internal/model"
	"github.com/waisee/microservices-go/inventory/internal/repository/converter"
	"github.com/waisee/microservices-go/inventory/internal/repository/record"
	"github.com/waisee/microservices-go/inventory/internal/service/input"
)

type PartRepository struct {
	mu    sync.RWMutex
	parts map[string]record.PartRecord
}

func NewPartRepository() *PartRepository {
	parts := map[string]record.PartRecord{
		"550e8400-e29b-41d4-a716-446655440001": {
			UUID:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Алюминиевый корпус",
			Description:   "Лёгкий корпус для небольших кораблей",
			Price:         500000, // 5000₽
			PartType:      record.PartTypeHull,
			StockQuantity: 10,
			CreatedAt:     time.Now(),
		},
		"550e8400-e29b-41d4-a716-446655440002": {
			UUID:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Титановый корпус",
			Description:   "Прочный корпус для средних кораблей",
			Price:         1500000, // 15000₽
			PartType:      record.PartTypeHull,
			StockQuantity: 5,
			CreatedAt:     time.Now(),
		},
		"550e8400-e29b-41d4-a716-446655440003": {
			UUID:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Ионный двигатель C",
			Description:   "Базовый ионный двигатель класса C",
			Price:         300000, // 3000₽
			PartType:      record.PartTypeEngine,
			StockQuantity: 8,
			CreatedAt:     time.Now(),
		},
		"550e8400-e29b-41d4-a716-446655440004": {
			UUID:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Ионный двигатель B",
			Description:   "Улучшенный ионный двигатель класса B",
			Price:         800000, // 8000₽
			PartType:      record.PartTypeEngine,
			StockQuantity: 3,
			CreatedAt:     time.Now(),
		},
		"550e8400-e29b-41d4-a716-446655440005": {
			UUID:          "550e8400-e29b-41d4-a716-446655440005",
			Name:          "Энергетический щит",
			Description:   "Стандартный энергетический щит",
			Price:         400000, // 4000₽
			PartType:      record.PartTypeShield,
			StockQuantity: 6,
			CreatedAt:     time.Now(),
		},
		"550e8400-e29b-41d4-a716-446655440006": {
			UUID:          "550e8400-e29b-41d4-a716-446655440006",
			Name:          "Лазерная пушка",
			Description:   "Точная лазерная пушка",
			Price:         250000, // 2500₽
			PartType:      record.PartTypeWeapon,
			StockQuantity: 7,
			CreatedAt:     time.Now(),
		},
		"550e8400-e29b-41d4-a716-446655440007": {
			UUID:          "550e8400-e29b-41d4-a716-446655440007",
			Name:          "Плазменный корпус",
			Description:   "Экспериментальный корпус (нет на складе)",
			Price:         2000000, // 20000₽
			PartType:      record.PartTypeHull,
			StockQuantity: 0,
			CreatedAt:     time.Now(),
		},
	}
	return &PartRepository{
		parts: parts,
	}
}

func (r *PartRepository) Get(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}
	return converter.RecordToModel(part), nil
}

func (r *PartRepository) List(ctx context.Context, filter input.PartFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(filter.UUIDs) > 0 {
		parts := make([]model.Part, 0)
		for _, uuid := range filter.UUIDs {
			part, ok := r.parts[uuid]
			if !ok {
				return nil, errs.ErrPartNotFound
			}
			parts = append(parts, converter.RecordToModel(part))
		}
		return parts, nil
	}

	if filter.PartType != input.PartTypeUnspecified {
		parts := make([]model.Part, 0)
		for _, part := range r.parts {
			if part.PartType == record.PartType(filter.PartType) {
				parts = append(parts, converter.RecordToModel(part))
			}
		}
		sort.Slice(parts, func(i, j int) bool {
			return parts[i].Name < parts[j].Name
		})
		return parts, nil
	}

	parts := make([]model.Part, 0)
	for _, part := range r.parts {
		parts = append(parts, converter.RecordToModel(part))
	}
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})
	return parts, nil
}

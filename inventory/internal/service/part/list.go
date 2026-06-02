package part

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/inventory/internal/errors"
	"github.com/waisee/microservices-go/inventory/internal/model"
	"github.com/waisee/microservices-go/inventory/internal/service/input"
)

func (s *PartService) List(ctx context.Context, filter input.PartFilter) ([]model.Part, error) {
	if len(filter.UUIDs) > 0 {
		for _, partUuid := range filter.UUIDs {
			if _, err := uuid.Parse(partUuid); err != nil {
				return nil, fmt.Errorf("получить список деталей: %w", errs.ErrInvalidUUID)
			}
		}
	}
	parts, err := s.repository.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	slog.InfoContext(ctx, "список деталей получен", "parts", parts)

	return parts, nil
}

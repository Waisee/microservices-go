package part

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/inventory/internal/errors"
	"github.com/waisee/microservices-go/inventory/internal/model"
)

func (s *PartService) Get(ctx context.Context, partUuid string) (model.Part, error) {
	if partUuid == "" {
		return model.Part{}, fmt.Errorf("получить деталь: %w", errs.ErrInvalidUUID)
	}
	if _, err := uuid.Parse(partUuid); err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", errs.ErrInvalidUUID)
	}
	part, err := s.repository.Get(ctx, partUuid)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	slog.InfoContext(ctx, "деталь получена", "part", part)
	return part, nil
}

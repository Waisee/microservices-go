package v1

import (
	"context"

	"github.com/waisee/microservices-go/inventory/internal/model"
	"github.com/waisee/microservices-go/inventory/internal/service/input"
)

type PartService interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, filter input.PartFilter) ([]model.Part, error)
}

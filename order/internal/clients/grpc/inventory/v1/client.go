package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/waisee/microservices-go/order/internal/clients/grpc/inventory/v1/converter"
	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/shared/pkg/maputil"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
)

const (
	timeout = 5 * time.Second
)

type InventoryClient struct {
	client inventoryv1.InventoryServiceClient
}

func NewInventoryClient(c inventoryv1.InventoryServiceClient) *InventoryClient {
	return &InventoryClient{
		client: c,
	}
}

func (c *InventoryClient) ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	uuidsStr := lo.Map(uuids, maputil.ToLoMap(func(id uuid.UUID) string {
		return id.String()
	}))

	resp, err := c.client.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuidsStr,
	})
	if err != nil {
		return nil, err
	}

	return lo.Map(resp.Parts, maputil.ToLoMap(converter.ProtoToModel)), nil
}

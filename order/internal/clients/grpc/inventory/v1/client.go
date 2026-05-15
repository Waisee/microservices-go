package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/clients/grpc/inventory/v1/converter"
	"github.com/waisee/microservices-go/order/internal/model"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
)

type InventoryClient struct {
	inventoryv1.InventoryServiceClient
}

func NewInventoryClient(grpc inventoryv1.InventoryServiceClient) *InventoryClient {
	return &InventoryClient{
		InventoryServiceClient: grpc,
	}
}

func (c *InventoryClient) ListParts(ctx context.Context, uuids []uuid.UUID) ([]model.Part, error) {
	uuidsStr := make([]string, 0, len(uuids))
	for _, uuid := range uuids {
		uuidsStr = append(uuidsStr, uuid.String())
	}
	resp, err := c.InventoryServiceClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuidsStr,
	})
	if err != nil {
		return nil, err
	}
	return converter.ProtoToModelList(resp.Parts), nil
}

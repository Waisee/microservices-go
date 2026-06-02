package v1

import inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"

type InventoryApi struct {
	partService PartService
	inventoryv1.UnimplementedInventoryServiceServer
}

func NewInventoryApi(partService PartService) *InventoryApi {
	return &InventoryApi{
		partService: partService,
	}
}

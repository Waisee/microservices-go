package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/waisee/microservices-go/inventory/internal/api/converter"
	"github.com/waisee/microservices-go/inventory/internal/errors"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, converter.ProtoToPartFilter(req))
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, errs.ErrInvalidUUID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.ModelToProtoList(parts),
	}, nil
}

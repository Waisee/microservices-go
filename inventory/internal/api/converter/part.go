package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/waisee/microservices-go/inventory/internal/model"
	"github.com/waisee/microservices-go/inventory/internal/service/input"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
)

func PartModelToProto(m model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          m.UUID,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		PartType:      PartTypeToProto(m.PartType),
		StockQuantity: m.StockQuantity,
		CreatedAt:     timestamppb.New(m.CreatedAt),
	}
}

func ModelToProto(m model.Part) *inventoryv1.GetPartResponse {
	return &inventoryv1.GetPartResponse{
		Part: PartModelToProto(m),
	}
}

func ProtoToModel(proto *inventoryv1.GetPartResponse) model.Part {
	return model.Part{
		UUID:          proto.Part.Uuid,
		Name:          proto.Part.Name,
		Description:   proto.Part.Description,
		Price:         proto.Part.Price,
		PartType:      PartTypeToModel(proto.Part.PartType),
		StockQuantity: proto.Part.StockQuantity,
		CreatedAt:     proto.Part.CreatedAt.AsTime(),
	}
}

func PartTypeToProto(partType model.PartType) inventoryv1.PartType {
	switch partType {
	case model.PartTypeUnspecified:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL
	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	}
	return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
}

func PartTypeToModel(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_UNSPECIFIED:
		return model.PartTypeUnspecified
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	}
	return model.PartTypeUnspecified
}

func ProtoToPartFilter(proto *inventoryv1.ListPartsRequest) input.PartFilter {
	return input.PartFilter{
		PartType: PartTypeToInput(proto.PartType),
		UUIDs:    proto.Uuids,
	}
}

func PartTypeToInput(partType inventoryv1.PartType) input.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_UNSPECIFIED:
		return input.PartTypeUnspecified
	case inventoryv1.PartType_PART_TYPE_HULL:
		return input.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return input.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return input.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return input.PartTypeWeapon
	}
	return input.PartTypeUnspecified
}

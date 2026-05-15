package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/waisee/microservices-go/order/internal/model"
	inventoryv1 "github.com/waisee/microservices-go/shared/pkg/proto/inventory/v1"
)

func ModelToProtoList(parts []model.Part) []*inventoryv1.Part {
	protos := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		protos = append(protos, ModelToProto(part))
	}
	return protos
}

func ModelToProto(part model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      partTypeToProto(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     timestamppb.New(part.CreatedAt),
	}
}

func partTypeToProto(partType model.PartType) inventoryv1.PartType {
	switch partType {
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

func ProtoToModelList(parts []*inventoryv1.Part) []model.Part {
	models := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		models = append(models, ProtoToModel(part))
	}
	return models
}

func ProtoToModel(part *inventoryv1.Part) model.Part {
	return model.Part{
		UUID:          uuid.MustParse(part.Uuid),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      partTypeToModel(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt.AsTime(),
	}
}

func partTypeToModel(partType inventoryv1.PartType) model.PartType {
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

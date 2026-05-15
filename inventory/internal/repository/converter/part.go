package converter

import (
	"github.com/waisee/microservices-go/inventory/internal/model"
	"github.com/waisee/microservices-go/inventory/internal/repository/record"
)

func RecordToModel(record record.PartRecord) model.Part {
	return model.Part{
		UUID:          record.UUID,
		Name:          record.Name,
		Description:   record.Description,
		Price:         record.Price,
		PartType:      model.PartType(record.PartType),
		StockQuantity: record.StockQuantity,
		CreatedAt:     record.CreatedAt,
	}
}

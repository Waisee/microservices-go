package converter

import (
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/order/internal/repository/record"
	"github.com/waisee/microservices-go/shared/pkg/maputil"
)

func OrderItemToRecord(item model.OrderItem) record.OrderItemRecord {
	return record.OrderItemRecord{
		PartUUID: item.PartUUID.String(),
		PartType: string(item.PartType),
		Price:    item.Price,
	}
}

func RecordToOrderItem(item record.OrderItemRecord) model.OrderItem {
	return model.OrderItem{
		PartUUID: uuid.MustParse(item.PartUUID),
		PartType: model.PartType(item.PartType),
		Price:    item.Price,
	}
}

func ModelToRecord(model model.Order) record.OrderRecord {
	transactionUUID := ""
	if model.TransactionUUID != nil {
		transactionUUID = model.TransactionUUID.String()
	}
	paymentMethod := ""
	if model.PaymentMethod != nil {
		paymentMethod = string(*model.PaymentMethod)
	}
	status := string(model.Status)
	return record.OrderRecord{
		UUID:            model.UUID.String(),
		Items:           lo.Map(model.Items, maputil.ToLoMap(OrderItemToRecord)),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          status,
		CreatedAt:       model.CreatedAt,
	}
}

func RecordToModel(record record.OrderRecord) model.Order {
	return model.Order{
		UUID:  uuid.MustParse(record.UUID),
		Items: lo.Map(record.Items, maputil.ToLoMap(RecordToOrderItem)),
		TransactionUUID: func() *uuid.UUID {
			if record.TransactionUUID == "" {
				return nil
			}
			result := uuid.MustParse(record.TransactionUUID)
			return &result
		}(),
		PaymentMethod: func() *model.PaymentMethod {
			if record.PaymentMethod == "" {
				return nil
			}
			result := model.PaymentMethod(record.PaymentMethod)
			return &result
		}(),
		Status:    model.OrderStatus(record.Status),
		CreatedAt: record.CreatedAt,
	}
}

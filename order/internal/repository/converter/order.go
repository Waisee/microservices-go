package converter

import (
	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/order/internal/repository/record"
)

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
		Items:           ModelToRecordItems(model.Items),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          status,
		CreatedAt:       model.CreatedAt,
	}
}

func ModelToRecordItems(items []model.OrderItem) []record.OrderItemRecord {
	result := make([]record.OrderItemRecord, 0, len(items))
	for _, item := range items {
		result = append(result, record.OrderItemRecord{
			PartUUID: item.PartUUID.String(),
			PartType: string(item.PartType),
			Price:    item.Price,
		})
	}
	return result
}

func RecordToModel(record record.OrderRecord) model.Order {
	return model.Order{
		UUID:  uuid.MustParse(record.UUID),
		Items: RecordToModelItems(record.Items),
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

func RecordToModelItems(items []record.OrderItemRecord) []model.OrderItem {
	result := make([]model.OrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, model.OrderItem{
			PartUUID: uuid.MustParse(item.PartUUID),
			PartType: model.PartType(item.PartType),
			Price:    item.Price,
		})
	}
	return result
}

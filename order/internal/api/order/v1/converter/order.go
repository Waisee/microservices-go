package converter

import (
	"github.com/google/uuid"

	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/order/internal/service/input"
	orderv1 "github.com/waisee/microservices-go/shared/pkg/openapi/order/v1"
)

func ProtoToCreateOrderInput(req *orderv1.CreateOrderRequest) input.CreateOrderInput {
	in := input.CreateOrderInput{
		HullUUID:   req.GetHullUUID(),
		EngineUUID: req.GetEngineUUID(),
	}

	if req.GetShieldUUID().IsSet() && !req.GetShieldUUID().IsNull() {
		v := req.GetShieldUUID().Value
		in.ShieldUUID = &v
	}
	if req.GetWeaponUUID().IsSet() && !req.GetWeaponUUID().IsNull() {
		v := req.GetWeaponUUID().Value
		in.WeaponUUID = &v
	}

	return in
}

func ModelToCreateOrderRes(order model.Order) (orderv1.CreateOrderRes, error) {
	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
	}, nil
}

func ModelToGetOrderRes(order model.Order) (orderv1.GetOrderRes, error) {
	hullUUID := uuid.Nil
	engineUUID := uuid.Nil
	shieldUUID := orderv1.OptNilUUID{
		Value: uuid.Nil,
		Set:   false,
		Null:  true,
	}
	weaponUUID := orderv1.OptNilUUID{
		Value: uuid.Nil,
		Set:   false,
		Null:  true,
	}
	for _, item := range order.Items {
		switch item.PartType {
		case model.PartTypeHull:
			hullUUID = item.PartUUID
		case model.PartTypeEngine:
			engineUUID = item.PartUUID
		case model.PartTypeShield:
			shieldUUID = orderv1.OptNilUUID{
				Value: item.PartUUID,
				Set:   true,
				Null:  false,
			}
		case model.PartTypeWeapon:
			weaponUUID = orderv1.OptNilUUID{
				Value: item.PartUUID,
				Set:   true,
				Null:  false,
			}
		}
	}
	transactionUUID := orderv1.OptNilUUID{
		Value: uuid.Nil,
		Set:   false,
		Null:  true,
	}

	if order.TransactionUUID != nil {
		transactionUUID = orderv1.OptNilUUID{
			Value: *order.TransactionUUID,
			Set:   true,
			Null:  false,
		}
	}
	paymentMethod := orderv1.OptNilPaymentMethod{
		Value: "",
		Set:   false,
		Null:  true,
	}
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.OptNilPaymentMethod{
			Value: orderv1.PaymentMethod(*order.PaymentMethod),
			Set:   true,
			Null:  false,
		}
	}
	status := orderv1.OrderStatus(order.Status)
	return &orderv1.OrderDto{
		HullUUID:        hullUUID,
		EngineUUID:      engineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		OrderUUID:       order.UUID,
		TotalPrice:      order.TotalPrice(),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          status,
		CreatedAt:       order.CreatedAt,
	}, nil
}

func ProtoToPaymentMethod(req *orderv1.PayOrderRequest) model.PaymentMethod {
	return model.PaymentMethod(req.GetPaymentMethod())
}

func ModelToPayOrderRes(transactionUUID uuid.UUID) (orderv1.PayOrderRes, error) {
	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

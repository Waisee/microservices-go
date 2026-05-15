package payment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/waisee/microservices-go/payment/internal/errors"
	"github.com/waisee/microservices-go/payment/internal/service/input"
)

func (s *PaymentService) Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error) {
	// 1. Валидация входных данных.
	if !in.PaymentMethod.IsValid() {
		return uuid.Nil, fmt.Errorf("pay order: %w", errs.ErrInvalidPaymentMethod)
	}
	// 2. Генерация transaction_uuid.
	transactionUUID := uuid.New()

	// 3. Логирование.
	slog.InfoContext(ctx, "оплата выполнена",
		"order_uuid", in.OrderUUID,
		"transaction_uuid", transactionUUID,
		"payment_method", in.PaymentMethod)

	return transactionUUID, nil
}

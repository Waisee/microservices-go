package payment_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/waisee/microservices-go/payment/internal/errors"
	"github.com/waisee/microservices-go/payment/internal/model"
	"github.com/waisee/microservices-go/payment/internal/service/input"
	"github.com/waisee/microservices-go/payment/internal/service/payment"
)

func TestPay(t *testing.T) {
	t.Parallel()

	type args struct {
		in input.PayOrderInput
	}
	type expected struct {
		err error
	}

	var (
		ctx       = context.Background()
		orderUUID = uuid.New()
	)

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "оплата картой",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "оплата через СБП",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodSBP,
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "оплата кредитной картой",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodCreditCard,
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "оплата деньгами инвестора",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodInvestorMoney,
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "PaymentMethod = UNSPECIFIED",
			args: args{
				in: input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodUnspecified,
				},
			},
			expected: expected{
				err: errs.ErrInvalidPaymentMethod,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			svc := payment.NewPaymentService()
			got, err := svc.Pay(ctx, test.args.in)
			if test.expected.err != nil {
				require.Error(t, err)
				assert.Equal(t, uuid.Nil, got)
				assert.ErrorIs(t, err, test.expected.err)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, got)
			}
		})
	}
}

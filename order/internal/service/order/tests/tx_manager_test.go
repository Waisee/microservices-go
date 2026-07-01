package order_test

import "context"

type fakeTxManager struct{}

func (fakeTxManager) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

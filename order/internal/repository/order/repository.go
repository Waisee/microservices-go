package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	squirrel "github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/waisee/microservices-go/order/internal/errors"
	"github.com/waisee/microservices-go/order/internal/model"
	"github.com/waisee/microservices-go/order/internal/repository/converter"
	"github.com/waisee/microservices-go/order/internal/repository/record"
)

// TxManager определяет контракт для управления транзакциями.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type OrderRepository struct {
	pool      *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	txManager TxManager
}

func NewOrderRepository(pool *pgxpool.Pool, txManager TxManager) *OrderRepository {
	return &OrderRepository{
		pool:      pool,
		getter:    trmpgx.DefaultCtxGetter,
		txManager: txManager,
	}
}

// Create атомарно сохраняет заказ и его строки в одной транзакции.
func (r *OrderRepository) Create(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.createOrder(txCtx, order); err != nil {
			return err
		}
		return r.createOrderItems(txCtx, order)
	})
}

// createOrder вставляет запись в таблицу orders.
func (r *OrderRepository) createOrder(ctx context.Context, order model.Order) error {
	rec := converter.ModelToRecord(order)

	const query = `
        INSERT INTO orders (uuid, status, created_at)
        VALUES ($1, $2, $3)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query,
		rec.UUID,
		rec.Status,
		rec.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}

// createOrderItems выполняет batch-insert строк заказа в order_items.
// order_uuid берётся из order.UUID (родитель агрегата), PartUUID/PartType/Price —
// из order.Items: в самой model.OrderItem ссылки на родителя нет.
func (r *OrderRepository) createOrderItems(ctx context.Context, order model.Order) error {
	if len(order.Items) == 0 {
		return nil
	}

	items := make([]record.OrderItemRecord, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, converter.OrderItemToRecord(item))
	}

	query := squirrel.Insert("order_items").
		Columns("order_uuid", "part_uuid", "part_type", "price").
		PlaceholderFormat(squirrel.Dollar)

	for _, item := range items {
		query = query.Values(order.UUID, item.PartUUID, item.PartType, item.Price)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("построить запрос: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("вставить строки заказа: %w", err)
	}

	return nil
}

func (r *OrderRepository) Get(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	const orderQuery = `
		SELECT uuid, status, transaction_uuid, payment_method, created_at
		FROM orders
		WHERE uuid = $1`

	var (
		order         model.Order
		paymentMethod *string
	)

	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, orderQuery, uuid).Scan(
		&order.UUID,
		&order.Status,
		&order.TransactionUUID,
		&paymentMethod,
		&order.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Order{}, errs.ErrOrderNotFound
	}
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	if paymentMethod != nil {
		pm := model.PaymentMethod(*paymentMethod)
		order.PaymentMethod = &pm
	}

	const itemsQuery = `
		SELECT part_uuid, part_type, price
		FROM order_items
		WHERE order_uuid = $1`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, itemsQuery, uuid)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить позиции заказа: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.PartUUID, &item.PartType, &item.Price); err != nil {
			return model.Order{}, fmt.Errorf("сканировать позицию заказа: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return model.Order{}, fmt.Errorf("прочитать позиции заказа: %w", err)
	}

	return order, nil
}

// Update обновляет заказ в транзакции — симметрично Create и готово
// к будущим мультистрочным обновлениям (например, синхронизация order_items).
func (r *OrderRepository) Update(ctx context.Context, order model.Order) error {
	return r.updateOrder(ctx, order)
}

func (r *OrderRepository) updateOrder(ctx context.Context, order model.Order) error {
	query := `UPDATE orders SET status = $1, transaction_uuid = $2, payment_method = $3, updated_at = $4 WHERE uuid = $5`

	rowsAffected, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, order.Status, order.TransactionUUID, order.PaymentMethod, time.Now(), order.UUID)
	if err != nil {
		return fmt.Errorf("обновить заказ: %w", err)
	}

	if rowsAffected.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}

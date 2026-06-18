package part

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/waisee/microservices-go/inventory/internal/errors"
	"github.com/waisee/microservices-go/inventory/internal/model"
	"github.com/waisee/microservices-go/inventory/internal/service/input"
)

// TxManager определяет контракт для управления транзакциями.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type PartRepository struct {
	pool      *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	txManager TxManager
}

func NewPartRepository(pool *pgxpool.Pool, txManager TxManager) *PartRepository {
	return &PartRepository{
		pool:      pool,
		getter:    trmpgx.DefaultCtxGetter,
		txManager: txManager,
	}
}

func (r *PartRepository) Get(ctx context.Context, uuid string) (model.Part, error) {
	query := `SELECT uuid, name, description, part_type, price, stock_quantity, created_at FROM parts WHERE uuid = $1`
	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, uuid)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Part{}, errs.ErrPartNotFound
	}
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}
	defer rows.Close()

	var part model.Part
	var found bool

	for rows.Next() {
		found = true
		err = rows.Scan(&part.UUID, &part.Name, &part.Description, &part.PartType, &part.Price, &part.StockQuantity, &part.CreatedAt)
		if err != nil {
			return model.Part{}, fmt.Errorf("сканировать деталь: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return model.Part{}, fmt.Errorf("прочитать деталь: %w", err)
	}

	if !found {
		return model.Part{}, errs.ErrPartNotFound
	}

	return part, nil
}

func (r *PartRepository) List(ctx context.Context, filter input.PartFilter) ([]model.Part, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.pool)

	if len(filter.UUIDs) > 0 {
		const query = `SELECT uuid, name, description, part_type, price, stock_quantity, created_at FROM parts WHERE uuid = ANY($1)`

		rows, err := db.Query(ctx, query, filter.UUIDs)
		if err != nil {
			return nil, fmt.Errorf("получить список деталей: %w", err)
		}
		defer rows.Close()

		found := make(map[string]model.Part, len(filter.UUIDs))
		for rows.Next() {
			var p model.Part
			if err := rows.Scan(&p.UUID, &p.Name, &p.Description, &p.PartType, &p.Price, &p.StockQuantity, &p.CreatedAt); err != nil {
				return nil, fmt.Errorf("сканировать деталь: %w", err)
			}
			found[p.UUID] = p
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("получить список деталей: %w", err)
		}

		parts := make([]model.Part, 0, len(filter.UUIDs))
		for _, u := range filter.UUIDs {
			p, ok := found[u]
			if !ok {
				return nil, errs.ErrPartNotFound
			}
			parts = append(parts, p)
		}
		return parts, nil
	}

	// Ветка по типу / все — сортировка по имени в SQL.
	var (
		query string
		args  []any
	)
	if filter.PartType != input.PartTypeUnspecified {
		query = `SELECT uuid, name, description, part_type, price, stock_quantity, created_at
		         FROM parts WHERE part_type = $1 ORDER BY name`
		args = []any{string(filter.PartType)}
	} else {
		query = `SELECT uuid, name, description, part_type, price, stock_quantity, created_at
		         FROM parts ORDER BY name`
	}
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	defer rows.Close()
	parts := make([]model.Part, 0)
	for rows.Next() {
		var p model.Part
		if err := rows.Scan(&p.UUID, &p.Name, &p.Description, &p.PartType, &p.Price, &p.StockQuantity, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("сканировать деталь: %w", err)
		}
		parts = append(parts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	return parts, nil
}

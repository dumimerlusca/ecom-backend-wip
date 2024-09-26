package model

import (
	"context"
	"ecom-backend/pkg/sqldb"
	"time"

	"github.com/lib/pq"
)

type ProductOptionRecord struct {
	Id        string
	ProductId string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type ProductOptionModel struct{}

func NewProductOptionModel() *ProductOptionModel {
	return &ProductOptionModel{}
}

func (p *ProductOptionModel) Insert(ctx context.Context, conn sqldb.Connection, option *ProductOptionRecord) (*ProductOptionRecord, error) {
	q := `INSERT INTO product_option (product_id, title) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, option.ProductId, option.Title).Scan(&option.Id, &option.CreatedAt, &option.UpdatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "duplicate_option_not_allowed"`:
			return nil, ErrDuplicatedProductOption
		default:
			return nil, err
		}
	}

	return option, nil
}

func (p *ProductOptionModel) DeleteAllByProductId(ctx context.Context, conn sqldb.Connection, productId string) error {
	q := `DELETE FROM product_option WHERE product_id = $1`
	_, err := conn.ExecContext(ctx, q, productId)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductOptionModel) FindForProducts(ctx context.Context, conn sqldb.Connection, productIds []string) (map[string][]*ProductOptionRecord, error) {
	q := `SELECT id, product_id, title, created_at, updated_at, deleted_at FROM product_option WHERE product_id = ANY($1)`

	rows, err := conn.QueryContext(ctx, q, pq.Array(productIds))

	if err != nil {
		return nil, err
	}

	resultMap := make(map[string][]*ProductOptionRecord)

	for rows.Next() {
		var record ProductOptionRecord

		err := rows.Scan(&record.Id, &record.ProductId, &record.Title, &record.CreatedAt, &record.UpdatedAt, &record.DeletedAt)

		if err != nil {
			return nil, err
		}

		if resultMap[record.ProductId] == nil {
			resultMap[record.ProductId] = []*ProductOptionRecord{}
		}

		resultMap[record.ProductId] = append(resultMap[record.ProductId], &record)
	}

	return resultMap, nil
}

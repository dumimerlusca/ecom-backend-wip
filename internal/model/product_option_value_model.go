package model

import (
	"context"
	"ecom-backend/pkg/sqldb"
	"time"
)

type ProductOptionValueRecord struct {
	Id        string
	OptionId  string
	VariantId string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type ProductOptionValueModel struct{}

func NewProductOptionValueModel() *ProductOptionValueModel {
	return &ProductOptionValueModel{}
}

func (m *ProductOptionValueModel) Insert(ctx context.Context, conn sqldb.Connection, value *ProductOptionValueRecord) (*ProductOptionValueRecord, error) {
	q := `INSERT INTO product_option_value (option_id, variant_id, title) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, value.OptionId, value.VariantId, value.Title).Scan(&value.Id, &value.CreatedAt, &value.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return value, nil
}

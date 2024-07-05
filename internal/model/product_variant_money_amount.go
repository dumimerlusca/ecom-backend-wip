package model

import (
	"context"
	"ecom-backend/pkg/sqldb"
	"time"
)

type ProductVariantMoneyAmountRecord struct {
	VariantId     string
	MoneyAmountId string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type ProductVariantMoneyAmountModel struct{}

func NewProductVariantMoneyAmountModel() *ProductVariantMoneyAmountModel {
	return &ProductVariantMoneyAmountModel{}
}

func (m *ProductVariantMoneyAmountModel) Insert(ctx context.Context, conn sqldb.Connection, value *ProductVariantMoneyAmountRecord) (*ProductVariantMoneyAmountRecord, error) {
	q := `INSERT INTO product_variant_money_amount (variant_id, money_amount_id) VALUES ($1, $2) RETURNING created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, value.VariantId, value.MoneyAmountId).Scan(&value.CreatedAt, &value.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return value, nil
}

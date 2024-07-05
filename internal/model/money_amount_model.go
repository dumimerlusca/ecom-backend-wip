package model

import (
	"context"
	"ecom-backend/pkg/sqldb"
	"time"
)

type MoneyAmountRecord struct {
	Id           string
	CurrencyCode string
	Amount       float32
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type MoneyAmountModel struct{}

func NewMoneyAmountModel() *MoneyAmountModel {
	return &MoneyAmountModel{}
}

func (m *MoneyAmountModel) Insert(ctx context.Context, conn sqldb.Connection, moneyAmount *MoneyAmountRecord) (*MoneyAmountRecord, error) {
	q := `INSERT INTO money_amount (currency_code, amount) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, moneyAmount.CurrencyCode, moneyAmount.Amount).Scan(&moneyAmount.Id, &moneyAmount.CreatedAt, &moneyAmount.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return moneyAmount, nil
}

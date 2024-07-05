package model

import (
	"context"
	"ecom-backend/pkg/sqldb"
	"time"
)

type ProductRecord struct {
	Id          string
	Title       string
	Subtitle    *string
	Description string
	Thumbnail   *string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type ProductModel struct {
}

func NewProductModel() *ProductModel {
	return &ProductModel{}
}

func (p *ProductModel) Insert(ctx context.Context, conn sqldb.Connection, product *ProductRecord) (*ProductRecord, error) {
	q := `INSERT INTO product (title, subtitle, description, thumbnail, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, product.Title, product.Subtitle, product.Description, product.Thumbnail, product.Status).Scan(&product.Id, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return product, nil
}

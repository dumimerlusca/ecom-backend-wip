package model

import (
	"context"
	"database/sql"
	"ecom-backend/pkg/sqldb"
	"errors"
	"time"
)

type ProductCategoryRecord struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	ParentId  *string   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductCategoryModel struct{}

func NewProductCategoryModel() *ProductCategoryModel {
	return &ProductCategoryModel{}
}

func (p *ProductCategoryModel) Insert(ctx context.Context, conn sqldb.Connection, category *ProductCategoryRecord) (*ProductCategoryRecord, error) {
	q := `INSERT INTO product_category (name, parent_id) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, category.Name, category.ParentId).Scan(&category.Id, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return category, nil
}

func (p *ProductCategoryModel) FindById(ctx context.Context, conn sqldb.Connection, id string) (*ProductCategoryRecord, error) {
	q := `SELECT id, name, parent_id, created_at, updated_at FROM product_category WHERE id = $1`

	category := &ProductCategoryRecord{}

	err := conn.QueryRowContext(ctx, q, id).Scan(&category.Id, &category.Name, &category.ParentId, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return category, nil
}

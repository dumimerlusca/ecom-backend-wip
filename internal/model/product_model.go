package model

import (
	"context"
	"database/sql"
	"ecom-backend/pkg/sqldb"
	"errors"
	"time"
)

type ProductRecord struct {
	Id          string
	Title       string
	Subtitle    *string
	Description string
	ThumbnailId *string
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
	q := `INSERT INTO product (title, subtitle, description, thumbnail_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, product.Title, product.Subtitle, product.Description, product.ThumbnailId, product.Status).Scan(&product.Id, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductModel) FindById(ctx context.Context, conn sqldb.Connection, id string) (*ProductRecord, error) {
	q := `SELECT id, title, subtitle, description, thumbnail_id, status, created_at, updated_at, deleted_at FROM product WHERE id = $1`

	product := &ProductRecord{}

	err := conn.QueryRowContext(ctx, q, id).Scan(&product.Id, &product.Title, &product.Subtitle, &product.Description, &product.ThumbnailId, &product.Status, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return product, nil
}

func (p *ProductModel) Update(ctx context.Context, conn sqldb.Connection, product *ProductRecord) (*ProductRecord, error) {
	q := `UPDATE product SET title = $1, subtitle = $2, description = $3, thumbnail_id = $4, status = $5, updated_at = $6 WHERE id = $7`

	product.UpdatedAt = time.Now()

	_, err := conn.ExecContext(ctx, q, product.Title, product.Subtitle, product.Description, product.ThumbnailId, product.Status, product.UpdatedAt, product.Id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return product, nil
}

func (p *ProductModel) MarkAsDeleted(ctx context.Context, conn sqldb.Connection, id string) error {
	q := `UPDATE product SET deleted_at = $1 WHERE id = $2`

	res, err := conn.ExecContext(ctx, q, time.Now(), id)

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrRecordNotFound
	}

	return nil

}

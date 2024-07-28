package model

import (
	"context"
	"database/sql"
	"ecom-backend/pkg/sqldb"
	"errors"
	"time"
)

type ProductCategoryRecord struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	ParentId  *string    `json:"parent_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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
	q := `SELECT id, name, parent_id, created_at, updated_at, deleted_at FROM product_category WHERE id = $1`

	category := &ProductCategoryRecord{}

	err := conn.QueryRowContext(ctx, q, id).Scan(&category.Id, &category.Name, &category.ParentId, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt)

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

func (p *ProductCategoryModel) FindAll(ctx context.Context, conn sqldb.Connection) ([]*ProductCategoryRecord, error) {
	q := `SELECT id, name, parent_id, created_at, updated_at FROM product_category WHERE deleted_at IS NULL`

	rows, err := conn.QueryContext(ctx, q)

	if err != nil {
		return nil, err
	}

	categories := []*ProductCategoryRecord{}

	for rows.Next() {
		var c ProductCategoryRecord

		err := rows.Scan(&c.Id, &c.Name, &c.ParentId, &c.CreatedAt, &c.UpdatedAt)

		if err != nil {
			return nil, err
		}

		categories = append(categories, &c)
	}

	return categories, nil
}

func (p *ProductCategoryModel) MarkAsDeleted(ctx context.Context, conn sqldb.Connection, id string) error {
	q := `UPDATE product_category SET deleted_at = $1 WHERE id = $2`

	res, err := conn.ExecContext(ctx, q, time.Now(), id)

	if err != nil {
		return err
	}

	rowsAff, _ := res.RowsAffected()

	if rowsAff == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (p *ProductCategoryModel) Update(ctx context.Context, conn sqldb.Connection, record *ProductCategoryRecord) (*ProductCategoryRecord, error) {
	q := `UPDATE product_category SET name = $1, parent_id = $2, updated_at = $3 WHERE id = $4`

	record.UpdatedAt = time.Now()

	res, err := conn.ExecContext(ctx, q, record.Name, record.ParentId, record.UpdatedAt, record.Id)

	if err != nil {
		return nil, err
	}

	if rowsAff, _ := res.RowsAffected(); rowsAff == 0 {
		return nil, ErrRecordNotFound
	}

	return record, nil

}

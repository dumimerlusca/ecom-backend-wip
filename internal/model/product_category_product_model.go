package model

import (
	"context"
	"ecom-backend/pkg/sqldb"
	"fmt"
	"strings"
)

type ProductCategoryProductRecord struct {
	CategoryId string
	ProductId  string
}

type ProductCategoryProductModel struct{}

func NewProductCategoryProductModel() *ProductCategoryProductModel {
	return &ProductCategoryProductModel{}
}

func (m *ProductCategoryProductModel) Insert(ctx context.Context, conn sqldb.Connection, record *ProductCategoryProductRecord) (*ProductCategoryProductRecord, error) {
	q := `INSERT INTO product_category_product (category_id, product_id) VALUES ($1, $2)`
	_, err := conn.ExecContext(ctx, q, record.CategoryId, record.ProductId)

	if err != nil {
		msg := err.Error()
		fmt.Println("MSG:", msg)
		switch {
		case strings.Contains(msg, "pq: invalid input syntax for type"):
			return nil, ErrInvalidProductCategory
		case strings.Contains(msg, `pq: insert or update on table "product_category_product" violates foreign key constraint "product_category_product_category_id_fkey"`):
			return nil, ErrProductCategoryNotFound
		case msg == `pq: duplicate key value violates unique constraint "product_category_product_pkey"`:
			return nil, ErrDuplicatedProductCategoryForProduct
		}
		return nil, err
	}
	return record, nil
}

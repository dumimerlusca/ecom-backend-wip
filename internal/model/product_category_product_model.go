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

func (m *ProductCategoryProductModel) DeleteAllByProductId(ctx context.Context, conn sqldb.Connection, productId string) error {
	q := `DELETE FROM product_category_product WHERE product_id = $1`
	_, err := conn.ExecContext(ctx, q, productId)

	if err != nil {
		return err
	}

	return nil
}

func (m *ProductCategoryProductModel) FindCategoriesForProduct(ctx context.Context, conn sqldb.Connection, productId string) ([]*ProductCategoryRecord, error) {
	q := `SELECT pc.id, pc.name, pc.parent_id, pc.created_at, pc.updated_at FROM product_category_product as pcp INNER JOIN product_category as pc ON pcp.category_id = pc.id WHERE product_id = $1`

	rows, err := conn.QueryContext(ctx, q, productId)

	if err != nil {
		return nil, err
	}

	list := []*ProductCategoryRecord{}

	for rows.Next() {
		var category ProductCategoryRecord

		err := rows.Scan(&category.Id, &category.Name, &category.ParentId, &category.CreatedAt, &category.UpdatedAt)

		if err != nil {
			return nil, err
		}

		list = append(list, &category)
	}

	return list, nil
}

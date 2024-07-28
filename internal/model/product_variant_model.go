package model

import (
	"context"
	"database/sql"
	"ecom-backend/pkg/sqldb"
	"errors"
	"time"
)

type ProductVariantRecord struct {
	Id                string
	ProductId         string
	Title             string
	Sku               *string
	Barcode           *int
	Material          *string
	Weight            *float32
	Length            *float32
	Width             *float32
	Height            *float32
	InventoryQuantity int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

type ProductVariantModel struct{}

func NewProductVariantModel() *ProductVariantModel {
	return &ProductVariantModel{}
}

func (p *ProductVariantModel) Insert(ctx context.Context, conn sqldb.Connection, variant *ProductVariantRecord) (*ProductVariantRecord, error) {
	q := `INSERT INTO product_variant (product_id, title, sku, barcode, material, weight, length, width,height, inventory_quantity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, updated_at`

	err := conn.QueryRowContext(ctx, q, variant.ProductId, variant.Title, variant.Sku, variant.Barcode, variant.Material, variant.Weight, variant.Length, variant.Width, variant.Height, variant.InventoryQuantity).Scan(&variant.Id, &variant.CreatedAt, &variant.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return variant, nil
}

func (p *ProductVariantModel) FindById(ctx context.Context, conn sqldb.Connection, id string) (*ProductVariantRecord, error) {
	q := `SELECT id, product_id, title, sku, barcode, material, weight, length, width, height, inventory_quantity, created_at, updated_at, deleted_at FROM product_variant WHERE id = $1`

	var variant ProductVariantRecord

	err := conn.QueryRowContext(ctx, q, id).Scan(&variant.Id, &variant.ProductId, &variant.Title, &variant.Sku, &variant.Barcode, &variant.Material, &variant.Weight, &variant.Length, &variant.Width, &variant.Height, &variant.InventoryQuantity, &variant.CreatedAt, &variant.UpdatedAt, &variant.DeletedAt)

	if err != nil {
		switch {

		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &variant, nil
}

func (p *ProductVariantModel) Update(ctx context.Context, conn sqldb.Connection, variant *ProductVariantRecord) (*ProductVariantRecord, error) {
	q := `UPDATE product_variant SET title = $1, sku = $2, barcode = $3, material = $4, weight = $5, length = $6, width = $7, height = $8, inventory_quantity = $9, updated_at = $10 WHERE id = $11`

	variant.UpdatedAt = time.Now()

	_, err := conn.ExecContext(ctx, q, variant.Title, variant.Sku, variant.Barcode, variant.Material, variant.Weight, variant.Length, variant.Width, variant.Height, variant.InventoryQuantity, variant.UpdatedAt, variant.Id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return variant, nil
}

func (p *ProductVariantModel) FindAllByProductId(ctx context.Context, conn sqldb.Connection, productId string) ([]*ProductVariantRecord, error) {
	q := `SELECT id, product_id, title, sku, barcode, material, weight, length, width, height, inventory_quantity, created_at, updated_at, deleted_at FROM product_variant WHERE product_id = $1`

	rows, err := conn.QueryContext(ctx, q, productId)

	variants := []*ProductVariantRecord{}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var variant ProductVariantRecord

		err := rows.Scan(&variant.Id, &variant.ProductId, &variant.Title, &variant.Sku, &variant.Barcode, &variant.Material, &variant.Weight, &variant.Length, &variant.Width, &variant.Height, &variant.InventoryQuantity, &variant.CreatedAt, &variant.UpdatedAt, &variant.DeletedAt)

		if err != nil {
			return nil, err
		}

		variants = append(variants, &variant)

	}

	return variants, nil
}

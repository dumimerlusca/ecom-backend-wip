package model

import (
	"context"
	"ecom-backend/pkg/sqldb"
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

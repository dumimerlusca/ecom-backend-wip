package service

import (
	"context"
	"database/sql"
	"ecom-backend/internal/consts"
	"ecom-backend/internal/model"
	"ecom-backend/internal/validator"
	"fmt"
)

type ProductService struct {
	db     *sql.DB
	models *model.Models
}

func NewProductService(db *sql.DB, models *model.Models) *ProductService {
	return &ProductService{db: db, models: models}
}

type CreateProductInput struct {
	Title       string   `json:"title"`
	Subtitle    *string  `json:"subtitle"`
	Description string   `json:"description"`
	Thumbnail   *string  `json:"thumbnail"`
	Material    *string  `json:"material"`
	Images      []string `json:"images"`
	Brand       *string  `json:"brand"`
	Status      string   `json:"status"`
	Categories  []struct {
		Id string `json:"id"`
	} `json:"categories"`
	Options []struct {
		Title string `json:"title"`
	} `json:"options"`
	Variants []VariantInput `json:"variants"`
}

type VariantInput struct {
	Title             string `json:"title"`
	Sku               string `json:"sku"`
	Barcode           int    `json:"barcode"`
	InventoryQuantity int    `json:"inventory_quantity"`
	Options           []struct {
		Value string `json:"value"`
	} `json:"options"`
	Prices []struct {
		Code   string  `json:"code"`
		Amount float32 `json:"amount"`
	} `json:"prices"`
}

func (input *CreateProductInput) Validate(v *validator.Validator) {
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(input.Description != "", "description", "must be provided")
	v.Check(validator.In(input.Status, consts.StatusDraft, consts.StatusPublished), "status", "invalid status")
	v.Check(len(input.Variants) > 0, "variants", "at least one variant must be provided")

	if len(input.Variants) > 0 {
		for _, variant := range input.Variants {
			variant.Validate(v)
		}
	}

}

func (input *VariantInput) Validate(v *validator.Validator) {
	v.Check(input.Title != "", "variant.title", "must be provided")
	v.Check(input.InventoryQuantity >= 0, "variant.inventory_quantity", "should not be negative")
	v.Check(len(input.Prices) != 0, "variant.prices", "at least one price must be provided")

	if len(input.Prices) > 0 {
		for _, price := range input.Prices {
			v.Check(price.Code != "", "variant.price_code", "must be provided")
			v.Check(price.Amount > 0, "variant.price_amount", "must be greater than zero")
		}
	}

	if len(input.Options) > 0 {
		for _, option := range input.Options {
			v.Check(option.Value != "", "variant.option_value", "must be provided")
		}

	}

}

func (s *ProductService) CreateProduct(ctx context.Context, input *CreateProductInput) (any, error) {

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	// create product record
	// the product record is the general description of the product,
	// the actual product entity containing the price that is used for purchase, wishlist, cart is the product_variant
	productRecord := &model.ProductRecord{Title: input.Title, Subtitle: input.Subtitle, Description: input.Description, Thumbnail: input.Thumbnail, Status: input.Status}

	product, err := s.models.ProductModel.Insert(ctx, tx, productRecord)

	if err != nil {
		return nil, fmt.Errorf("failed to create product record: %w", err)
	}

	productOptionRecords := []*model.ProductOptionRecord{}

	// create product_option records
	// ex: size, color, etc
	for _, option := range input.Options {
		optionRecord, err := s.models.ProductOptionModel.Insert(ctx, tx, &model.ProductOptionRecord{ProductId: product.Id, Title: option.Title})

		if err != nil {
			return nil, fmt.Errorf("failed to create product_option record: %w", err)
		}

		productOptionRecords = append(productOptionRecords, optionRecord)
	}

	// create product_variant records
	// the entity that contains the price, inventory quantity, sku, barcode, etc and it's used in the cart, wishlist, purchase
	for _, variant := range input.Variants {
		productVariantRecord := &model.ProductVariantRecord{ProductId: product.Id, Title: variant.Title, Sku: &variant.Sku, Barcode: &variant.Barcode, InventoryQuantity: variant.InventoryQuantity}

		variantRecord, err := s.models.ProductVariantModel.Insert(ctx, tx, productVariantRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to create product_variant record: %w", err)
		}

		for i, option := range variant.Options {
			// creating product_option_value value record
			// each variant is linked to multiple product options
			// ex: variant 1 is linked to size: M, color: red
			_, err := s.models.ProductOptionValueModel.Insert(ctx, tx, &model.ProductOptionValueRecord{VariantId: variantRecord.Id, Title: option.Value, OptionId: productOptionRecords[i].Id})

			if err != nil {
				return nil, fmt.Errorf("failed to create product_option_value record: %w", err)
			}

		}

		for _, price := range variant.Prices {
			// create money_amount record
			// the price is stored in the money_amount table and will be linked to the variant using the product_variant_money_amount table
			moneyAmountRecord, err := s.models.MoneyAmountModel.Insert(ctx, tx, &model.MoneyAmountRecord{CurrencyCode: price.Code, Amount: price.Amount})

			if err != nil {
				return nil, fmt.Errorf("failed to create money_amount record: %w", err)
			}

			// create product_variant_money_amount record
			// the product_variant_money_amount table is used to link the product_variant to the money_amount
			// 1 -> M relationship, in order to support multiple currencies
			_, err = s.models.ProductVariantMoneyAmountModel.Insert(ctx, tx, &model.ProductVariantMoneyAmountRecord{VariantId: variantRecord.Id, MoneyAmountId: moneyAmountRecord.Id})

			if err != nil {
				return nil, fmt.Errorf("failed to create product_variant_money_amount record: %w", err)
			}

		}

	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return product, nil
}

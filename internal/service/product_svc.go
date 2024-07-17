package service

import (
	"context"
	"database/sql"
	"ecom-backend/internal/model"
	"errors"
	"fmt"
)

type ProductService struct {
	db     *sql.DB
	models *model.Models
}

func NewProductService(db *sql.DB, models *model.Models) *ProductService {
	return &ProductService{db: db, models: models}
}

func (svc *ProductService) CreateProduct(ctx context.Context, input *CreateProductInput) (*DetailedProduct, error) {

	productCategoryRecords := []*model.ProductCategoryRecord{}
	productOptionRecords := []*model.ProductOptionRecord{}
	variantRecords := []*model.ProductVariantRecord{}
	moneyAmountRecords := map[string][]*model.MoneyAmountRecord{}
	variantOptionValueRecords := map[string][]*model.ProductOptionValueRecord{}

	tx, err := svc.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	// create product record
	// the product record is the general description of the product,
	// the actual product entity containing the price that is used for purchase, wishlist, cart is the product_variant
	productRecord := &model.ProductRecord{Title: input.Title, Subtitle: input.Subtitle, Description: input.Description, ThumbnailId: input.ThumbnailId, Status: input.Status}

	product, err := svc.models.ProductModel.Insert(ctx, tx, productRecord)

	if err != nil {
		return nil, fmt.Errorf("failed to create product record: %w", err)
	}

	// create product_category_product records
	// the product_category_product table is used to link the product to multiple categories
	for _, category := range input.Categories {
		categoryRecord, err := svc.models.ProductCategoryModel.FindById(ctx, tx, category.Id)

		if err != nil {
			if errors.Is(err, model.ErrRecordNotFound) {
				return nil, model.ErrProductCategoryNotFound
			} else {
				return nil, err
			}
		}

		productCategoryRecords = append(productCategoryRecords, categoryRecord)

		_, err = svc.models.ProductCategoryProductModel.Insert(ctx, tx, &model.ProductCategoryProductRecord{ProductId: product.Id, CategoryId: category.Id})

		if err != nil {
			return nil, fmt.Errorf("failed to create product_category_product record: %w", err)
		}

	}

	// create product_option records
	// ex: size, color, etc
	for _, option := range input.Options {
		optionRecord, err := svc.models.ProductOptionModel.Insert(ctx, tx, &model.ProductOptionRecord{ProductId: product.Id, Title: option.Title})

		if err != nil {
			return nil, fmt.Errorf("failed to create product_option record: %w", err)
		}

		productOptionRecords = append(productOptionRecords, optionRecord)
	}

	// create product_variant records
	// the entity that contains the price, inventory quantity, sku, barcode, etc and it's used in the cart, wishlist, purchase
	for _, variant := range input.Variants {
		productVariantRecord := &model.ProductVariantRecord{ProductId: product.Id, Title: variant.Title, Sku: &variant.Sku, Barcode: &variant.Barcode, InventoryQuantity: variant.InventoryQuantity}

		variantRecord, err := svc.models.ProductVariantModel.Insert(ctx, tx, productVariantRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to create product_variant record: %w", err)
		}

		variantRecords = append(variantRecords, variantRecord)

		for i, option := range variant.Options {
			// creating product_option_value value record
			// each variant is linked to multiple product options
			// ex: variant 1 is linked to size: M, color: red
			productOptionValueRecord, err := svc.models.ProductOptionValueModel.Insert(ctx, tx, &model.ProductOptionValueRecord{VariantId: variantRecord.Id, Title: option.Value, OptionId: productOptionRecords[i].Id})

			if err != nil {
				return nil, fmt.Errorf("failed to create product_option_value record: %w", err)
			}

			variantOptionValueRecords[variantRecord.Id] = append(variantOptionValueRecords[variantRecord.Id], productOptionValueRecord)
		}

		for _, price := range variant.Prices {
			// create money_amount record
			// the price is stored in the money_amount table and will be linked to the variant using the product_variant_money_amount table
			moneyAmountRecord, err := svc.models.MoneyAmountModel.Insert(ctx, tx, &model.MoneyAmountRecord{CurrencyCode: price.Code, Amount: price.Amount})

			if err != nil {
				return nil, fmt.Errorf("failed to create money_amount record: %w", err)
			}

			// create product_variant_money_amount record
			// the product_variant_money_amount table is used to link the product_variant to the money_amount
			// 1 -> M relationship, in order to support multiple currencies
			_, err = svc.models.ProductVariantMoneyAmountModel.Insert(ctx, tx, &model.ProductVariantMoneyAmountRecord{VariantId: variantRecord.Id, MoneyAmountId: moneyAmountRecord.Id})

			if err != nil {
				return nil, fmt.Errorf("failed to create product_variant_money_amount record: %w", err)
			}

			moneyAmountRecords[variantRecord.Id] = append(moneyAmountRecords[variantRecord.Id], moneyAmountRecord)
		}

	}

	// link product to images
	for _, image := range input.Images {
		_, err := svc.models.EntityFileModel.Insert(ctx, tx, &model.EntityFileRecord{EntityId: product.Id, FileId: image.Id})

		if err != nil {
			return nil, fmt.Errorf("failed to link image to product: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// populate the final product
	finalProduct := BuildDetailedProduct(product, productCategoryRecords, productOptionRecords, variantRecords, moneyAmountRecords, variantOptionValueRecords)

	return finalProduct, nil
}

func (svc *ProductService) UpdateProduct(ctx context.Context, productId string, input *UpdateProductInput) (*model.ProductRecord, error) {
	tx, err := svc.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	productRecord, err := svc.models.ProductModel.FindById(ctx, tx, productId)

	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		productRecord.Title = *input.Title
	}

	if input.Subtitle != nil {
		productRecord.Subtitle = input.Subtitle
	}

	if input.Description != nil {
		productRecord.Description = *input.Description
	}

	if input.Status != nil {
		productRecord.Status = *input.Status
	}

	if input.ThumbnailId != nil {
		productRecord.ThumbnailId = input.ThumbnailId
	}

	productRecord, err = svc.models.ProductModel.Update(ctx, tx, productRecord)

	if err != nil {
		return nil, err
	}

	if input.Categories != nil {
		// handle categories

		// delete all existing product_category_product records
		// and create new ones
		err := svc.models.ProductCategoryProductModel.DeleteAllByProductId(ctx, tx, productId)

		if err != nil {
			return nil, err

		}

		for _, category := range *input.Categories {
			_, err := svc.models.ProductCategoryProductModel.Insert(ctx, tx, &model.ProductCategoryProductRecord{ProductId: productId, CategoryId: category.Id})

			if err != nil {
				return nil, fmt.Errorf("failed to create product_category_product record: %w", err)
			}
		}
	}

	if input.Options != nil {
		// handle product options

		// delete all existing product_option records
		// and create new ones
		err := svc.models.ProductOptionModel.DeleteAllByProductId(ctx, tx, productId)

		if err != nil {
			return nil, err
		}

		for _, option := range *input.Options {
			_, err := svc.models.ProductOptionModel.Insert(ctx, tx, &model.ProductOptionRecord{ProductId: productId, Title: option.Title})

			if err != nil {
				return nil, fmt.Errorf("failed to create product_option record: %w", err)
			}

		}

	}

	if input.Images != nil {
		// handle images

		// delete all existing entity_file records
		// and create new ones
		err := svc.models.EntityFileModel.DeleteAllByEntityId(ctx, tx, productId)

		if err != nil {
			return nil, err
		}

		for _, image := range *input.Images {
			_, err := svc.models.EntityFileModel.Insert(ctx, tx, &model.EntityFileRecord{EntityId: productId, FileId: image.Id})

			if err != nil {
				return nil, fmt.Errorf("failed to link image to product: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return productRecord, nil
}

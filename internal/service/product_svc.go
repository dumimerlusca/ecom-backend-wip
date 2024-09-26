package service

import (
	"context"
	"database/sql"
	"ecom-backend/internal/model"
	"ecom-backend/pkg/sqldb"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type ProductService struct {
	db     *sql.DB
	models *model.Models
}

func NewProductService(db *sql.DB, models *model.Models) *ProductService {
	return &ProductService{db: db, models: models}
}

func (svc *ProductService) CreateProduct(ctx context.Context, input *CreateProductInput) (*AggregateProduct, error) {

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

	imageIds := []string{}

	// link product to images
	for _, image := range input.Images {
		imageIds = append(imageIds, image.Id)

		_, err := svc.models.EntityFileModel.Insert(ctx, tx, &model.EntityFileRecord{EntityId: product.Id, FileId: image.Id})

		if err != nil {
			return nil, fmt.Errorf("failed to link image to product: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	aggFields := BuildAggregateFieldsList(imageIds, productCategoryRecords, productOptionRecords, variantRecords, moneyAmountRecords, variantOptionValueRecords)

	aggProduct := BuildAggregateProduct(product, aggFields)

	return aggProduct, nil
}

func (svc *ProductService) UpdateProductDetails(ctx context.Context, productId string, input *UpdateProductInput) (*model.ProductRecord, error) {
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

func (svc *ProductService) UpdateVariantDetails(ctx context.Context, variantId string, input *UpdateVariantInput) (any, error) {
	tx, err := svc.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	variantRecord, err := svc.models.ProductVariantModel.FindById(ctx, tx, variantId)

	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		variantRecord.Title = *input.Title
	}

	if input.Sku != nil {
		variantRecord.Sku = input.Sku
	}

	if input.Barcode != nil {
		variantRecord.Barcode = input.Barcode
	}

	if input.InventoryQuantity != nil {
		variantRecord.InventoryQuantity = *input.InventoryQuantity
	}

	variantRecord, err = svc.models.ProductVariantModel.Update(ctx, tx, variantRecord)

	if err != nil {
		return nil, err
	}

	if input.Options != nil {
		// handle variant option values

		// delete existing product_option_value records
		// and create new ones
		err := svc.models.ProductOptionValueModel.DeleteAllByVariantId(ctx, tx, variantId)

		if err != nil {
			return nil, err
		}

		for _, option := range *input.Options {
			_, err := svc.models.ProductOptionValueModel.Insert(ctx, tx, &model.ProductOptionValueRecord{VariantId: variantId, Title: option.Value, OptionId: option.Id})

			if err != nil {
				return nil, fmt.Errorf("failed to create product_option_value record: %w", err)
			}
		}
	}

	if input.Prices != nil {
		// handle prices

		// delete existing product_variant_money_amount records
		err := svc.models.ProductVariantMoneyAmountModel.DeleteAllByVariantId(ctx, tx, variantId)

		if err != nil {
			return nil, fmt.Errorf("error while deleting product_variant_money_amount records")
		}

		moneyAmountRecords := []*model.MoneyAmountRecord{}

		// create new money amount records
		for _, price := range *input.Prices {
			moneyAmountRecord, err := svc.models.MoneyAmountModel.Insert(ctx, tx, &model.MoneyAmountRecord{CurrencyCode: price.Code, Amount: price.Amount})

			if err != nil {
				return nil, fmt.Errorf("error while creating money_amount records")
			}

			moneyAmountRecords = append(moneyAmountRecords, moneyAmountRecord)

		}

		// create new product_variant_money_amount records
		for _, moneyAmount := range moneyAmountRecords {
			svc.models.ProductVariantMoneyAmountModel.Insert(ctx, tx, &model.ProductVariantMoneyAmountRecord{VariantId: variantId, MoneyAmountId: moneyAmount.Id})
		}

	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return variantRecord, nil

}

func (svc *ProductService) ListProducts(ctx context.Context, opt ProductListingOptions) ([]*model.ProductRecord, int, error) {
	limit := opt.PageSize
	offset := (opt.Page - 1) * opt.PageSize

	totalCountQ := `SELECT COUNT(*) from product WHERE deleted_at IS NULL`
	var totalCount int
	err := svc.db.QueryRowContext(ctx, totalCountQ).Scan(&totalCount)

	if err != nil {
		return nil, 0, err
	}

	q := `SELECT id, title, subtitle, description, thumbnail_id, status, created_at, updated_at, deleted_at FROM product WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := svc.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	productsList := []*model.ProductRecord{}

	for rows.Next() {
		var product model.ProductRecord

		err := rows.Scan(&product.Id, &product.Title, &product.Subtitle, &product.Description, &product.ThumbnailId, &product.Status, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)

		if err != nil {
			return nil, 0, err
		}

		productsList = append(productsList, &product)
	}

	return productsList, totalCount, nil
}

type ProductListingOptions struct {
	Page     uint
	PageSize uint
}

func (svc *ProductService) ListAggregateProducts(ctx context.Context, opt ProductListingOptions) ([]*AggregateProduct, int, error) {
	products, rowCount, err := svc.ListProducts(ctx, opt)

	if err != nil {
		return nil, 0, err
	}

	productIds := []string{}

	for _, p := range products {
		productIds = append(productIds, p.Id)
	}

	aggFieldsMap, err := svc.GetAggregateFieldsForProductsList(ctx, productIds)

	if err != nil {
		return nil, 0, err
	}

	aggProductList := []*AggregateProduct{}

	for _, p := range products {
		aggProductList = append(aggProductList, BuildAggregateProduct(p, aggFieldsMap[p.Id]))
	}

	return aggProductList, rowCount, nil
}

func (svc *ProductService) GetAggregateFieldsForProductsList(ctx context.Context, productIds []string) (map[string]*AggregateProductListFields, error) {
	variantsMap, err := svc.models.ProductVariantModel.FindAllByProductIds(ctx, svc.db, productIds)

	if err != nil {
		return nil, err
	}

	categoriesMap, err := svc.models.ProductCategoryProductModel.FindCategoriesForProducts(ctx, svc.db, productIds)

	if err != nil {
		return nil, err
	}

	productOptionsMap, err := svc.models.ProductOptionModel.FindForProducts(ctx, svc.db, productIds)

	if err != nil {
		return nil, err
	}

	variantPricesMap, err := svc.getVariantPricesMap(ctx, svc.db, productIds)

	if err != nil {
		return nil, err
	}

	variantOptionValuesMap, err := svc.getVariantOptionValuesMap(ctx, svc.db, productIds)

	if err != nil {
		return nil, err
	}

	imageIds, err := svc.models.EntityFileModel.FindAllFilesByEntityId(ctx, svc.db, productIds)

	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]*AggregateProductListFields)

	for _, id := range productIds {
		resultMap[id] = BuildAggregateFieldsList(imageIds[id], categoriesMap[id], productOptionsMap[id], variantsMap[id], variantPricesMap[id], variantOptionValuesMap[id])
	}

	return resultMap, nil
}

func (svc *ProductService) GetAggregateProductById(ctx context.Context, id string) (*AggregateProduct, error) {
	product, err := svc.models.ProductModel.FindById(ctx, svc.db, id)

	if err != nil {
		return nil, err
	}

	aggFieldsMap, err := svc.GetAggregateFieldsForProductsList(ctx, []string{id})

	if err != nil {
		return nil, err
	}

	return BuildAggregateProduct(product, aggFieldsMap[id]), nil
}

func (svc *ProductService) getVariantPricesMap(ctx context.Context, conn sqldb.Connection, productIds []string) (map[string]map[string][]*model.MoneyAmountRecord, error) {
	q := `SELECT ma.id, ma.currency_code, ma.amount, ma.created_at, ma.updated_at, pvma.variant_id, pv.product_id
		  FROM product_variant_money_amount AS pvma 
		  INNER JOIN money_amount AS ma ON pvma.money_amount_id = ma.id
		  INNER JOIN product_variant AS pv ON pvma.variant_id = pv.id
		  WHERE pv.product_id = ANY($1)`

	rows, err := conn.QueryContext(ctx, q, pq.Array(productIds))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pricesMap := map[string]map[string][]*model.MoneyAmountRecord{}

	for rows.Next() {
		var record model.MoneyAmountRecord
		var variantId string
		var productId string

		err := rows.Scan(&record.Id, &record.CurrencyCode, &record.Amount, &record.CreatedAt, &record.UpdatedAt, &variantId, &productId)

		if err != nil {
			return nil, err
		}

		if pricesMap[productId] == nil {
			pricesMap[productId] = make(map[string][]*model.MoneyAmountRecord)
		}

		if pricesMap[productId][variantId] == nil {
			pricesMap[productId][variantId] = []*model.MoneyAmountRecord{}
		}

		pricesMap[productId][variantId] = append(pricesMap[productId][variantId], &record)

	}

	return pricesMap, nil
}

func (svc *ProductService) getVariantOptionValuesMap(ctx context.Context, conn sqldb.Connection, productIds []string) (map[string]map[string][]*model.ProductOptionValueRecord, error) {
	q := `SELECT pov.id, pov.option_id, pov.variant_id, pov.title, pov.created_at, pov.updated_at, pov.deleted_at, pv.product_id FROM product_option_value AS pov
		  INNER JOIN product_variant as pv ON pv.id = pov.variant_id
		  WHERE pv.product_id = ANY($1)`

	rows, err := conn.QueryContext(ctx, q, pq.Array(productIds))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	optionValuesMap := make(map[string]map[string][]*model.ProductOptionValueRecord)

	for rows.Next() {
		var record model.ProductOptionValueRecord
		var productId string

		err := rows.Scan(&record.Id, &record.OptionId, &record.VariantId, &record.Title, &record.CreatedAt, &record.UpdatedAt, &record.DeletedAt, &productId)

		if err != nil {
			return nil, err
		}

		variantId := record.VariantId

		if optionValuesMap[productId] == nil {
			optionValuesMap[productId] = make(map[string][]*model.ProductOptionValueRecord)
		}

		if optionValuesMap[productId][variantId] == nil {
			optionValuesMap[productId][variantId] = []*model.ProductOptionValueRecord{}
		}

		optionValuesMap[productId][variantId] = append(optionValuesMap[productId][variantId], &record)

	}

	return optionValuesMap, nil
}

func (svc *ProductService) MarkProductAsDeleted(ctx context.Context, productId string) error {
	err := svc.models.ProductModel.MarkAsDeleted(ctx, svc.db, productId)

	if err != nil {
		return err
	}

	return nil
}

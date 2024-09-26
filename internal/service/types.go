package service

import (
	"ecom-backend/internal/consts"
	"ecom-backend/internal/model"
	"ecom-backend/internal/validator"
	"time"
)

type AggregateProductListFields struct {
	Variants   []AggregateProductVariant `json:"variants"`
	Categories []ProductCategoryInfo     `json:"categories"`
	Options    []ProductOptionDTO        `json:"options"`
	Images     []ProductImage            `json:"images"`
}

type AggregateProduct struct {
	Id          string                    `json:"id"`
	Title       string                    `json:"title"`
	Subtitle    *string                   `json:"subtitle"`
	Description string                    `json:"description"`
	Thumbnail   *ProductImage             `json:"thumbnail"`
	Status      string                    `json:"status"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
	DeletedAt   *time.Time                `json:"deleted_at"`
	Variants    []AggregateProductVariant `json:"variants"`
	Categories  []ProductCategoryInfo     `json:"categories"`
	Options     []ProductOptionDTO        `json:"options"`
	Images      []ProductImage            `json:"images"`
}

type AggregateProductVariant struct {
	Id                string               `json:"id"`
	ProductId         string               `json:"product_id"`
	Title             string               `json:"title"`
	Sku               *string              `json:"sku"`
	Barcode           *int                 `json:"barcode"`
	Material          *string              `json:"material"`
	Weight            *float32             `json:"weight"`
	Length            *float32             `json:"length"`
	Width             *float32             `json:"width"`
	Height            *float32             `json:"height"`
	InventoryQuantity int                  `json:"inventory_quantity"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	DeletedAt         *time.Time           `json:"deleted_at"`
	Prices            []VariantPriceDTO    `json:"prices"`
	Options           []VariantOptionValue `json:"options"`
}

type ProductOptionDTO struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
type ProductImage struct {
	Id string `json:"id"`
}

type ProductCategoryInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type VariantPriceDTO struct {
	Id           string     `json:"id"`
	CurrencyCode string     `json:"currency_code"`
	Amount       float32    `json:"amount"`
	CreateAt     time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type VariantOptionValue struct {
	Id       string `json:"id"`
	Value    string `json:"value"`
	OptionId string `json:"option_id"`
}

func BuildAggregateProduct(productRecord *model.ProductRecord, aggListFields *AggregateProductListFields) *AggregateProduct {
	p := AggregateProduct{}

	p.Id = productRecord.Id
	p.Title = productRecord.Title
	p.Subtitle = productRecord.Subtitle
	p.Description = productRecord.Description
	p.Status = productRecord.Status
	p.CreatedAt = productRecord.CreatedAt
	p.UpdatedAt = productRecord.UpdatedAt
	p.DeletedAt = productRecord.DeletedAt

	if productRecord.ThumbnailId != nil {
		p.Thumbnail = &ProductImage{Id: *productRecord.ThumbnailId}

	}

	p.Variants = aggListFields.Variants
	p.Options = aggListFields.Options
	p.Categories = aggListFields.Categories
	p.Images = aggListFields.Images

	return &p
}

func BuildAggregateFieldsList(
	productImages []string,
	categoryRecords []*model.ProductCategoryRecord,
	optionRecords []*model.ProductOptionRecord,
	variantRecords []*model.ProductVariantRecord,
	variantMoneyAmountRecords map[string][]*model.MoneyAmountRecord,
	variantOptionValueRecords map[string][]*model.ProductOptionValueRecord) *AggregateProductListFields {

	agg := AggregateProductListFields{}

	for _, option := range optionRecords {
		agg.Options = append(agg.Options, ProductOptionDTO{Id: option.Id, Title: option.Title})
	}

	for _, imageId := range productImages {
		agg.Images = append(agg.Images, ProductImage{Id: imageId})
	}

	agg.Categories = []ProductCategoryInfo{}

	for _, categoryRecord := range categoryRecords {
		pci := ProductCategoryInfo{}
		pci.Id = categoryRecord.Id
		pci.Name = categoryRecord.Name

		agg.Categories = append(agg.Categories, pci)
	}

	agg.Variants = []AggregateProductVariant{}

	for _, variantRecord := range variantRecords {
		dpv := AggregateProductVariant{}
		dpv.Id = variantRecord.Id
		dpv.ProductId = variantRecord.ProductId
		dpv.Title = variantRecord.Title
		dpv.Sku = variantRecord.Sku
		dpv.Barcode = variantRecord.Barcode
		dpv.Material = variantRecord.Material
		dpv.Weight = variantRecord.Weight
		dpv.Length = variantRecord.Length
		dpv.Width = variantRecord.Width
		dpv.Height = variantRecord.Height
		dpv.InventoryQuantity = variantRecord.InventoryQuantity
		dpv.CreatedAt = variantRecord.CreatedAt
		dpv.UpdatedAt = variantRecord.UpdatedAt
		dpv.DeletedAt = variantRecord.DeletedAt

		dpv.Prices = []VariantPriceDTO{}
		for _, moneyAmountRecord := range variantMoneyAmountRecords[variantRecord.Id] {
			vp := VariantPriceDTO{}
			vp.Id = moneyAmountRecord.Id
			vp.CurrencyCode = moneyAmountRecord.CurrencyCode
			vp.Amount = moneyAmountRecord.Amount
			vp.CreateAt = moneyAmountRecord.CreatedAt
			vp.UpdatedAt = moneyAmountRecord.UpdatedAt
			vp.DeletedAt = moneyAmountRecord.DeletedAt

			dpv.Prices = append(dpv.Prices, vp)
		}

		dpv.Options = []VariantOptionValue{}

		for _, optionValueRecord := range variantOptionValueRecords[variantRecord.Id] {
			vov := VariantOptionValue{}
			vov.Id = optionValueRecord.Id
			vov.Value = optionValueRecord.Title
			vov.OptionId = optionValueRecord.OptionId

			dpv.Options = append(dpv.Options, vov)
		}

		agg.Variants = append(agg.Variants, dpv)
	}

	return &agg
}

type CreateProductInput struct {
	Title       string              `json:"title"`
	Subtitle    *string             `json:"subtitle"`
	Description string              `json:"description"`
	ThumbnailId *string             `json:"thumbnail_id"`
	Material    *string             `json:"material"`
	Images      []ProductImageInput `json:"images"`
	Brand       *string             `json:"brand"`
	Status      string              `json:"status"`
	Categories  []struct {
		Id string `json:"id"`
	} `json:"categories"`
	Options  []ProductOptionInput        `json:"options"`
	Variants []CreateProductVariantInput `json:"variants"`
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

type CreateProductVariantInput struct {
	Title             string `json:"title"`
	Sku               string `json:"sku"`
	Barcode           int    `json:"barcode"`
	InventoryQuantity int    `json:"inventory_quantity"`
	Options           []struct {
		Value string `json:"value"`
	} `json:"options"`
	Prices []PriceInput `json:"prices"`
}

type PriceInput struct {
	Code   string  `json:"code"`
	Amount float32 `json:"amount"`
}

func (input *CreateProductVariantInput) Validate(v *validator.Validator) {
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

type ProductImageInput struct {
	Id string `json:"id"`
}

type ProductCategoryInput struct {
	Id string `json:"id"`
}
type ProductOptionInput struct {
	Title string `json:"title"`
}

type UpdateProductInput struct {
	Title       *string                 `json:"title"`
	Subtitle    *string                 `json:"subtitle"`
	Description *string                 `json:"description"`
	ThumbnailId *string                 `json:"thumbnail_id"`
	Material    *string                 `json:"material"`
	Images      *[]ProductImageInput    `json:"images"`
	Brand       *string                 `json:"brand"`
	Status      *string                 `json:"status"`
	Categories  *[]ProductCategoryInput `json:"categories"`
	Options     *[]ProductOptionInput   `json:"options"`
}

func (input *UpdateProductInput) Validate(v *validator.Validator) {
	if input.Title != nil {
		v.Check(*input.Title != "", "title", "must not be empty")
	}

	if input.Description != nil {
		v.Check(*input.Description != "", "description", "must not be empty")
	}

	if input.Status != nil {
		v.Check(validator.In(*input.Status, consts.StatusDraft, consts.StatusPublished), "status", "invalid status")
	}
}

type UpdateVariantInput struct {
	Title             *string `json:"title"`
	Sku               *string `json:"sku"`
	Barcode           *int    `json:"barcode"`
	InventoryQuantity *int    `json:"inventory_quantity"`
	Options           *[]struct {
		Value string `json:"value"`
		Id    string `json:"id"`
	} `json:"options"`
	Prices *[]PriceInput `json:"prices"`
}

func (input *UpdateVariantInput) Validate(v *validator.Validator) {
	if input.Title != nil {
		v.Check(*input.Title != "", "title", "must not be empty")
	}

	if input.InventoryQuantity != nil {
		v.Check(*input.InventoryQuantity >= 0, "inventory_quantity", "should not be negative")
	}

	if input.Prices != nil {
		for _, price := range *input.Prices {
			v.Check(price.Code != "", "price.code", "must not be empty")
			v.Check(price.Amount > 0, "price.amount", "must be greater than zero")
		}
	}

	if input.Options != nil {
		for _, option := range *input.Options {
			v.Check(option.Value != "", "option.value", "must not be empty")
			v.Check(option.Id != "", "option.id", "must not be empty")
		}
	}

}

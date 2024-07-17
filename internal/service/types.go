package service

import (
	"ecom-backend/internal/consts"
	"ecom-backend/internal/model"
	"ecom-backend/internal/validator"
	"time"
)

type DetailedProduct struct {
	Id          string                   `json:"id"`
	Title       string                   `json:"title"`
	Subtitle    *string                  `json:"subtitle"`
	Description string                   `json:"description"`
	Thumbnail   *ProductImage            `json:"thumbnail"`
	Status      string                   `json:"status"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	DeletedAt   *time.Time               `json:"deleted_at"`
	Images      []ProductImage           `json:"images"`
	Variants    []DetailedProductVariant `json:"variants"`
	Categories  []ProductCategoryInfo    `json:"categories"`
}

type DetailedProductVariant struct {
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
	Prices            []VariantPrice       `json:"prices"`
	Options           []VariantOptionValue `json:"options"`
}

type ProductImage struct {
	Id string `json:"id"`
}

type ProductCategoryInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type VariantPrice struct {
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

func BuildDetailedProduct(
	productRecord *model.ProductRecord,
	categoryRecords []*model.ProductCategoryRecord,
	optionRecords []*model.ProductOptionRecord,
	variantRecords []*model.ProductVariantRecord,
	variantMoneyAmountRecords map[string][]*model.MoneyAmountRecord,
	variantOptionValueRecords map[string][]*model.ProductOptionValueRecord) *DetailedProduct {

	dp := DetailedProduct{}

	dp.Id = productRecord.Id
	dp.Title = productRecord.Title
	dp.Subtitle = productRecord.Subtitle
	dp.Description = productRecord.Description
	dp.Status = productRecord.Status
	dp.CreatedAt = productRecord.CreatedAt
	dp.UpdatedAt = productRecord.UpdatedAt
	dp.DeletedAt = productRecord.DeletedAt

	if productRecord.ThumbnailId != nil {
		dp.Thumbnail = &ProductImage{Id: *productRecord.ThumbnailId}

	}

	dp.Categories = []ProductCategoryInfo{}

	for _, categoryRecord := range categoryRecords {
		pci := ProductCategoryInfo{}
		pci.Id = categoryRecord.Id
		pci.Name = categoryRecord.Name

		dp.Categories = append(dp.Categories, pci)
	}

	dp.Variants = []DetailedProductVariant{}

	for _, variantRecord := range variantRecords {
		dpv := DetailedProductVariant{}
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

		dpv.Prices = []VariantPrice{}
		for _, moneyAmountRecord := range variantMoneyAmountRecords[variantRecord.Id] {
			vp := VariantPrice{}
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

		dp.Variants = append(dp.Variants, dpv)
	}

	return &dp
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
	Options  []ProductOptionInput `json:"options"`
	Variants []VariantInput       `json:"variants"`
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

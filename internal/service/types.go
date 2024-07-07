package service

import (
	"ecom-backend/internal/model"
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

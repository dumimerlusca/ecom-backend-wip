package model

import "ecom-backend/pkg/sqldb"

type Models struct {
	ProductModel                   *ProductModel
	ProductVariantModel            *ProductVariantModel
	ProductOptionModel             *ProductOptionModel
	ProductOptionValueModel        *ProductOptionValueModel
	MoneyAmountModel               *MoneyAmountModel
	ProductVariantMoneyAmountModel *ProductVariantMoneyAmountModel
	ProductCategoryModel           *ProductCategoryModel
	ProductCategoryProductModel    *ProductCategoryProductModel
	FileModel                      *FileModel
	EntityFileModel                *EntityFileModel
}

func NewModels(conn sqldb.Connection) *Models {
	return &Models{
		ProductModel:                   NewProductModel(),
		ProductVariantModel:            NewProductVariantModel(),
		ProductOptionModel:             NewProductOptionModel(),
		ProductOptionValueModel:        NewProductOptionValueModel(),
		MoneyAmountModel:               NewMoneyAmountModel(),
		ProductVariantMoneyAmountModel: NewProductVariantMoneyAmountModel(),
		ProductCategoryModel:           NewProductCategoryModel(),
		ProductCategoryProductModel:    NewProductCategoryProductModel(),
		FileModel:                      NewFileModel(),
		EntityFileModel:                NewEntityFileModel(),
	}
}

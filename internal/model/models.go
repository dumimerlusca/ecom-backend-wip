package model

import "ecom-backend/pkg/sqldb"

type Models struct {
	ProductModel                   *ProductModel
	ProductVariantModel            *ProductVariantModel
	ProductOptionModel             *ProductOptionModel
	ProductOptionValueModel        *ProductOptionValueModel
	MoneyAmountModel               *MoneyAmountModel
	ProductVariantMoneyAmountModel *ProductVariantMoneyAmountModel
}

func NewModels(conn sqldb.Connection) *Models {
	return &Models{
		ProductModel:                   NewProductModel(),
		ProductVariantModel:            NewProductVariantModel(),
		ProductOptionModel:             NewProductOptionModel(),
		ProductOptionValueModel:        NewProductOptionValueModel(),
		MoneyAmountModel:               NewMoneyAmountModel(),
		ProductVariantMoneyAmountModel: NewProductVariantMoneyAmountModel(),
	}
}

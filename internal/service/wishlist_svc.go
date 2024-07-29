package service

import (
	"context"
	"database/sql"
	"ecom-backend/internal/model"

	"github.com/lib/pq"
)

type WishlistService struct {
	db            *sql.DB
	wishlistModel *model.WishlistModel
}

func NewWishlistService(db *sql.DB, wishlistModel *model.WishlistModel) *WishlistService {
	return &WishlistService{db: db, wishlistModel: wishlistModel}
}

func (svc *WishlistService) Insert(ctx context.Context, userIdentifier string, variantId string) (*model.WishlistRecord, error) {
	record, err := svc.wishlistModel.Insert(ctx, svc.db, &model.WishlistRecord{UserIdentifier: userIdentifier, VariantId: variantId})

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (svc *WishlistService) DeleteItem(ctx context.Context, userIdentifier string, recordId string) error {
	wishlistItem, err := svc.wishlistModel.FindById(ctx, svc.db, recordId)

	if err != nil {
		return err
	}

	if wishlistItem.UserIdentifier != userIdentifier {
		return ErrUnauthorizedRequest
	}

	err = svc.wishlistModel.DeleteById(ctx, svc.db, recordId)

	return err
}

func (svc *WishlistService) ListProducts(ctx context.Context, userIdentifier string) ([]*WishlistItemDTO, error) {
	// get product and variant info
	q := `SELECT p.id, p.title, p.subtitle, p.description, p.thumbnail_id, 
			pv.id, pv.title FROM wishlist AS w
			LEFT JOIN product_variant AS pv ON w.variant_id = pv.id
			LEFT JOIN product AS p ON p.id = pv.product_id
			WHERE user_identifier = $1`

	rows, err := svc.db.QueryContext(ctx, q, userIdentifier)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	wishlistItems := []*WishlistItemDTO{}
	variantIds := []string{}

	for rows.Next() {
		item := WishlistItemDTO{Prices: []VariantPriceDTO{}, Options: []VariantOptionValueDTO{}}

		err := rows.Scan(&item.ProductId, &item.ProductTitle, &item.ProductSubtitle, &item.ProductDescription, &item.ThumbnailId, &item.VariantId, &item.VariantTitle)

		if err != nil {
			return nil, err
		}

		wishlistItems = append(wishlistItems, &item)
		variantIds = append(variantIds, item.VariantId)
	}

	pricesMap := map[string][]VariantPriceDTO{}

	// get variant prices
	q = `SELECT ma.id, ma.currency_code, ma.amount, pvma.variant_id FROM money_amount AS ma
		INNER JOIN product_variant_money_amount AS pvma ON pvma.money_amount_id = ma.id
		WHERE variant_id = ANY($1)`

	rows, err = svc.db.QueryContext(ctx, q, pq.Array(variantIds))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var priceItem VariantPriceDTO
		var variantId string

		err := rows.Scan(&priceItem.Id, &priceItem.CurrencyCode, &priceItem.Amount, &variantId)

		if err != nil {
			return nil, err
		}

		pricesMap[variantId] = append(pricesMap[variantId], priceItem)
	}

	optionValuesMap := map[string][]VariantOptionValueDTO{}

	// get variant option values
	q = `SELECT pov.variant_id, pov.id, pov.title, po.id, po.title FROM product_option_value AS pov
		 LEFT JOIN product_option AS po ON po.id = pov.option_id
		 WHERE variant_id = ANY($1)`

	rows, err = svc.db.QueryContext(ctx, q, pq.Array(variantIds))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var optionValueItem VariantOptionValueDTO
		var variantId string

		err := rows.Scan(&variantId, &optionValueItem.Id, &optionValueItem.Value, &optionValueItem.OptionId, &optionValueItem.OptionTitle)

		if err != nil {
			return nil, err
		}

		optionValuesMap[variantId] = append(optionValuesMap[variantId], optionValueItem)
	}

	for i := 0; i < len(wishlistItems); i++ {
		wishlistItems[i].Prices = pricesMap[wishlistItems[i].VariantId]
		wishlistItems[i].Options = optionValuesMap[wishlistItems[i].VariantId]
	}

	return wishlistItems, nil
}

type WishlistItemDTO struct {
	ProductId          string                  `json:"product_id"`
	ProductTitle       string                  `json:"product_title"`
	ProductSubtitle    *string                 `json:"product_subtitle"`
	ProductDescription string                  `json:"product_description"`
	ThumbnailId        *string                 `json:"thumbnail_id"`
	VariantId          string                  `json:"variant_id"`
	VariantTitle       string                  `json:"variant_title"`
	Prices             []VariantPriceDTO       `json:"prices"`
	Options            []VariantOptionValueDTO `json:"options"`
}

type VariantOptionValueDTO struct {
	Id          string `json:"id"`
	Value       string `json:"value"`
	OptionId    string `json:"option_id"`
	OptionTitle string `json:"option_title"`
}

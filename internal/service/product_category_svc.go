package service

import (
	"context"
	"database/sql"
	"ecom-backend/internal/model"
	"errors"
)

type ProductCategoryService struct {
	models *model.Models
	db     *sql.DB
}

func NewProductCategoryService(db *sql.DB, models *model.Models) *ProductCategoryService {
	return &ProductCategoryService{db: db, models: models}
}

func (svc *ProductCategoryService) CreateProductCategory(ctx context.Context, name string, parentId *string) (*model.ProductCategoryRecord, error) {
	if parentId != nil && *parentId != "" {
		// check if parent category exists
		_, err := svc.models.ProductCategoryModel.FindById(ctx, svc.db, *parentId)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrRecordNotFound):
				return nil, model.ErrParentProductCategoryNotFound
			default:
				return nil, err
			}
		}

	}

	// create product category record
	productCategoryRecord, err := svc.models.ProductCategoryModel.Insert(ctx, svc.db, &model.ProductCategoryRecord{Name: name, ParentId: parentId})

	return productCategoryRecord, err
}

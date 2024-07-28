package service

import (
	"context"
	"database/sql"
	"ecom-backend/internal/model"
	"ecom-backend/internal/validator"
	"errors"
	"time"
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

type ProductCategoryWithChildren struct {
	Id        string                         `json:"id"`
	Name      string                         `json:"name"`
	Children  []*ProductCategoryWithChildren `json:"children"`
	CreatedAt time.Time                      `json:"created_at"`
	UpdatedAt time.Time                      `json:"updated_at"`
}

func (svc *ProductCategoryService) GetAll(ctx context.Context) ([]*ProductCategoryWithChildren, error) {
	categories, err := svc.models.ProductCategoryModel.FindAll(ctx, svc.db)

	if err != nil {
		return nil, err
	}

	categoryMap := map[string]*ProductCategoryWithChildren{}

	for _, record := range categories {
		categoryMap[record.Id] = &ProductCategoryWithChildren{Id: record.Id,
			Name:      record.Name,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			Children:  []*ProductCategoryWithChildren{},
		}
	}

	nestedList := []*ProductCategoryWithChildren{}

	for _, record := range categories {
		if record.ParentId == nil {
			nestedList = append(nestedList, categoryMap[record.Id])
		} else {
			categoryMap[*record.ParentId].Children = append(categoryMap[*record.ParentId].Children, categoryMap[record.Id])
		}
	}

	return nestedList, nil
}

func (svc *ProductCategoryService) MarkAsDeleted(ctx context.Context, categpryId string) error {
	err := svc.models.ProductCategoryModel.MarkAsDeleted(ctx, svc.db, categpryId)

	return err
}

type UpdateProductCategoryInput struct {
	Name     *string `json:"name"`
	ParentId *string `json:"parent_id"`
}

func (input *UpdateProductCategoryInput) Validate(v *validator.Validator) {
	if input.Name != nil {
		v.Check(*input.Name != "", "name", "must not be empty")
	}
}

func (svc *ProductCategoryService) UpdateById(ctx context.Context, categoryId string, input UpdateProductCategoryInput) (*model.ProductCategoryRecord, error) {
	tx, err := svc.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	record, err := svc.models.ProductCategoryModel.FindById(ctx, tx, categoryId)

	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		record.Name = *input.Name
	}

	if input.ParentId != nil {
		if *input.ParentId == "" {
			record.ParentId = nil
		} else {
			// check if parent category exists
			_, err := svc.models.ProductCategoryModel.FindById(ctx, tx, *input.ParentId)

			if err != nil {
				return nil, model.ErrInvalidProductCategory
			}

			record.ParentId = input.ParentId
		}
	}

	record, err = svc.models.ProductCategoryModel.Update(ctx, tx, record)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return record, nil
}

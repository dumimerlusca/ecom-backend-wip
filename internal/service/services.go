package service

import (
	"database/sql"
	"ecom-backend/internal/model"
)

type Services struct {
	Product         *ProductService
	ProductCategory *ProductCategoryService
	Upload          *UploadService
}

func NewServices(db *sql.DB, models *model.Models) *Services {
	return &Services{
		Product:         NewProductService(db, models),
		ProductCategory: NewProductCategoryService(db, models),
		Upload:          NewUploadService(db, models.FileModel),
	}
}

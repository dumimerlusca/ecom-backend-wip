package service

import (
	"database/sql"
	"ecom-backend/internal/model"
)

type Services struct {
	Product         *ProductService
	ProductCategory *ProductCategoryService
	Upload          *UploadService
	Auth            *AuthService
	Token           *TokenService
}

func NewServices(db *sql.DB, models *model.Models) *Services {
	tokenSvc := NewTokenService(db, models.TokenModel, models.UserModel)

	return &Services{
		Product:         NewProductService(db, models),
		ProductCategory: NewProductCategoryService(db, models),
		Upload:          NewUploadService(db, models.FileModel),
		Token:           tokenSvc,
		Auth:            NewAuthService(db, models.UserModel, models.TokenModel, tokenSvc),
	}
}

package main

import "ecom-backend/internal/handlers"

type Handlers struct {
	product           *handlers.ProductHandler
	productCategories *handlers.ProductCategoryHandler
	fileUpload        *handlers.UploadHandler
	auth              *handlers.AuthHandler
}

func (app *application) createHandlers() *Handlers {
	return &Handlers{
		product:           handlers.NewProductHandler(app.logger, app.services.Product),
		productCategories: handlers.NewProductCategoryHandler(app.logger, app.services.ProductCategory),
		fileUpload:        handlers.NewUploadHandler(app.logger, app.services.Upload),
		auth:              handlers.NewAuthHandler(app.logger, app.services.Auth),
	}
}

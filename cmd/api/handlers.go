package main

import "ecom-backend/internal/handlers"

type Handlers struct {
	products handlers.ProductsHandler
}

func (app *application) createHandlers() *Handlers {
	return &Handlers{products: *handlers.NewProductsHandler(app.logger, app.services.Product)}
}

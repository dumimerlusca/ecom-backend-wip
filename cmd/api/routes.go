package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	h := app.createHandlers()

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Route not found 123"))
	})

	// SERVING ADMIN APP
	router.NotFound = http.FileServer(http.Dir("admin"))

	// API ROUTES
	router.HandlerFunc(http.MethodGet, "/api/v1/products", h.product.GetProducts)
	router.HandlerFunc(http.MethodPost, "/api/v1/products", h.product.CreateProduct)
	router.PATCH("/api/v1/products/:productId", h.product.UpdateProduct)
	router.HandlerFunc(http.MethodGet, "/api/v1/products/:productId", h.product.GetProduct)
	router.HandlerFunc(http.MethodDelete, "/api/v1/products/:productId", h.product.DeleteProduct)

	router.HandlerFunc(http.MethodPost, "/api/v1/product-categories", h.productCategories.Create)

	// File upload
	router.HandlerFunc(http.MethodPost, "/api/v1/upload", h.fileUpload.UploadFile)
	router.GET("/api/v1/files/:fileId", h.fileUpload.ServerFile)

	m := app.middleware

	return m.RecoverPanic(router)
}

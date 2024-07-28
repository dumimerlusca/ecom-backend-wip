package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	h := app.createHandlers()
	m := app.middleware

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Route not found 123"))
	})

	// Serving admin app
	router.NotFound = http.FileServer(http.Dir("admin"))

	// Api routes
	router.HandlerFunc(http.MethodGet, "/api/v1/products", h.product.GetProducts)
	router.HandlerFunc(http.MethodPost, "/api/v1/products", h.product.CreateProduct)
	router.PATCH("/api/v1/products/:productId", h.product.UpdateProductGeneralInfo)
	router.PATCH("/api/v1/variants/:variantId", h.product.UpdateVariantDetails)
	router.GET("/api/v1/products/:productId", h.product.GetProduct)
	router.DELETE("/api/v1/products/:productId", h.product.DeleteProduct)
	router.HandlerFunc(http.MethodPost, "/api/v1/product-categories", h.productCategories.Create)
	router.HandlerFunc(http.MethodGet, "/api/v1/product-categories", h.productCategories.GetAll)
	router.DELETE("/api/v1/product-categories/:categoryId", h.productCategories.DeleteById)
	router.PATCH("/api/v1/product-categories/:categoryId", h.productCategories.UpdateById)

	// File upload
	router.HandlerFunc(http.MethodPost, "/api/v1/upload", h.fileUpload.UploadFile)
	router.GET("/api/v1/files/:fileId", h.fileUpload.ServerFile)

	// Auth
	router.HandlerFunc(http.MethodPost, "/api/v1/users", h.auth.RegisterUser)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", h.auth.Login)

	return m.RecoverPanic(m.Authenticate(router))
}

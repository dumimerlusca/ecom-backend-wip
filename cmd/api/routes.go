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
	router.HandlerFunc(http.MethodGet, "/api/v1/products", h.products.GetProducts)
	router.HandlerFunc(http.MethodPost, "/api/v1/products", h.products.CreateProduct)
	router.HandlerFunc(http.MethodGet, "/api/v1/products/:productId", h.products.GetProduct)
	router.HandlerFunc(http.MethodPatch, "/api/v1/products/:productId", h.products.UpdateProduct)
	router.HandlerFunc(http.MethodDelete, "/api/v1/products/:productId", h.products.DeleteProduct)

	m := app.middleware

	return m.RecoverPanic(router)
}

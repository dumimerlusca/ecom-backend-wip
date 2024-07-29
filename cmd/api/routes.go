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
		w.Write([]byte("Route not found"))
	})

	// Serving admin app
	router.NotFound = http.FileServer(http.Dir("admin"))

	// Admin only routes
	router.POST("/api/v1/products", m.AdminOnly(h.product.CreateProduct))
	router.PATCH("/api/v1/products/:productId", m.AdminOnly(h.product.UpdateProductGeneralInfo))
	router.PATCH("/api/v1/variants/:variantId", m.AdminOnly(h.product.UpdateVariantDetails))
	router.DELETE("/api/v1/products/:productId", m.AdminOnly(h.product.DeleteProduct))
	router.POST("/api/v1/product-categories", m.AdminOnly(h.productCategories.Create))
	router.DELETE("/api/v1/product-categories/:categoryId", m.AdminOnly(h.productCategories.DeleteById))
	router.PATCH("/api/v1/product-categories/:categoryId", m.AdminOnly(h.productCategories.UpdateById))

	// Public routes
	router.GET("/api/v1/products", h.product.GetProducts)
	router.GET("/api/v1/products/:productId", h.product.GetProduct)
	router.GET("/api/v1/product-categories", h.productCategories.GetAll)
	router.POST("/api/v1/wishlist/add", m.RequireSessionOrUser(h.wishlist.Create))
	router.GET("/api/v1/wishlist", m.RequireSessionOrUser(h.wishlist.GetAll))
	router.DELETE("/api/v1/wishlist/remove/:id", m.RequireSessionOrUser(h.wishlist.DeleteItem))

	// File upload
	router.POST("/api/v1/upload", m.RequireActivation(h.fileUpload.UploadFile))
	router.GET("/api/v1/files/:fileId", m.RequireActivation(h.fileUpload.ServerFile))

	// Auth
	router.HandlerFunc(http.MethodPost, "/api/v1/users", h.auth.RegisterUser)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", h.auth.Login)
	router.GET("/api/v1/session", h.auth.GetSession)

	return m.RecoverPanic(m.EnableCORS(m.Authenticate(router)))
}

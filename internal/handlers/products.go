package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"ecom-backend/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ProductsHandler struct {
	BaseHandler
	productSrv *service.ProductService
}

func NewProductsHandler(logger *jsonlog.Logger, productSrv *service.ProductService) *ProductsHandler {
	return &ProductsHandler{BaseHandler: BaseHandler{logger: logger}, productSrv: productSrv}
}

func (h *ProductsHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	input := service.CreateProductInput{}

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Validate(v); !v.Valid() {
		h.FailedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = h.productSrv.CreateProduct(r.Context(), &input)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicateBarcode):
			h.BadRequestResponse(w, r, err)
		case errors.Is(err, model.ErrDuplicatedProductOption):
			h.BadRequestResponse(w, r, err)
		default:
			h.ServerErrorResponse(w, r, err)

		}

		return
	}

	w.Write([]byte("Product created"))
}

func (h *ProductsHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	h.WriteJson(w, 200, "Get products", nil)
}
func (h *ProductsHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetProduct")
}
func (h *ProductsHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateProduct")
}
func (h *ProductsHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteProduct")
}

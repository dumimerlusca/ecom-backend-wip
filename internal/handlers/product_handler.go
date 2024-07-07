package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"ecom-backend/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

type ProductHandler struct {
	BaseHandler
	productSvc *service.ProductService
}

func NewProductHandler(logger *jsonlog.Logger, productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{BaseHandler: BaseHandler{logger: logger}, productSvc: productSvc}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	input := service.CreateProductInput{}

	err := h.ReadJSON(w, r, &input)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Validate(v); !v.Valid() {
		h.FailedValidationResponse(w, r, v.Errors)
		return
	}

	detailedProduct, err := h.productSvc.CreateProduct(r.Context(), &input)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrProductCategoryNotFound),
			errors.Is(err, model.ErrDuplicatedProductOption),
			errors.Is(err, model.ErrDuplicatedProductCategoryForProduct),
			errors.Is(err, model.ErrInvalidProductCategory),
			errors.Is(err, model.ErrFileNotFound):
			h.BadRequestResponse(w, r, err)
		default:
			h.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = h.WriteJson(w, http.StatusCreated, ResponseBody{Payload: Envelope{"product": detailedProduct}}, nil)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	h.WriteJson(w, 200, "Get products", nil)
}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetProduct")
}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateProduct")
}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteProduct")
}

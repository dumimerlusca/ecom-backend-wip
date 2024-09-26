package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"ecom-backend/internal/validator"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ProductHandler struct {
	BaseHandler
	productSvc *service.ProductService
}

func NewProductHandler(logger *jsonlog.Logger, productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{BaseHandler: BaseHandler{logger: logger}, productSvc: productSvc}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (h *ProductHandler) UpdateProductGeneralInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	productId := ps.ByName("productId")

	if productId == "" {
		h.BadRequestResponse(w, r, errors.New("missing product id"))
		return
	}

	var input service.UpdateProductInput

	err := h.ReadJSON(w, r, &input)

	if err != nil {
		h.BadRequestResponse(w, r, err)
		return

	}

	v := validator.New()

	input.Validate(v)

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.Errors)
		return

	}

	product, err := h.productSvc.UpdateProductDetails(r.Context(), productId, &input)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			h.NotFoundResponse(w, r)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return

	}

	// TODO: change response to include detailed product
	h.WriteJson(w, http.StatusOK, ResponseBody{Payload: Envelope{"product_info": product}}, nil)
}

func (h *ProductHandler) UpdateVariantDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	variantId := ps.ByName("variantId")

	if variantId == "" {
		h.BadRequestResponse(w, r, errors.New("missing variant id"))
		return
	}

	var input service.UpdateVariantInput

	err := h.ReadJSON(w, r, &input)

	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	input.Validate(v)

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = h.productSvc.UpdateVariantDetails(r.Context(), variantId, &input)

	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			h.NotFoundResponse(w, r)
		} else {
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	h.WriteJson(w, http.StatusOK, Envelope{"success": true}, nil)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)

	if err != nil {
		h.BadRequestResponse(w, r, errors.New("invalid query parameter `page`"))
		return
	}

	if page <= 0 {
		h.BadRequestResponse(w, r, errors.New("page must be > 0"))
		return
	}

	pageSize, err := strconv.ParseInt(r.URL.Query().Get("pageSize"), 10, 64)

	if err != nil {
		h.BadRequestResponse(w, r, errors.New("invalid query parameter `pageSize`"))
		return
	}

	if pageSize <= 0 {
		h.BadRequestResponse(w, r, errors.New("pageSize must be > 0"))
		return
	}

	list, rowCount, err := h.productSvc.ListAggregateProducts(r.Context(), service.ProductListingOptions{Page: uint(page), PageSize: uint(pageSize)})

	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	h.WriteJson(w, http.StatusOK, ResponseBody{Payload: list, Metadata: PaginationMetadata{Page: int(page), PageSize: int(pageSize), RowsTotal: rowCount}}, nil)
}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	productId := ps.ByName("productId")

	if productId == "" {
		h.NotFoundResponse(w, r)
		return
	}

	product, err := h.productSvc.GetAggregateProductById(r.Context(), productId)

	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			h.NotFoundResponse(w, r)
			return
		}

		h.ServerErrorResponse(w, r, err)
		return
	}

	h.WriteJson(w, http.StatusOK, ResponseBody{Payload: Envelope{"product": product}}, nil)
}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	productId := ps.ByName("productId")

	if productId == "" {
		h.NotFoundResponse(w, r)
		return
	}

	err := h.productSvc.MarkProductAsDeleted(r.Context(), productId)

	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			h.NotFoundResponse(w, r)
		} else {
			h.ServerErrorResponse(w, r, err)
		}

		return
	}

	h.WriteJson(w, http.StatusOK, Envelope{"success": true}, nil)
}

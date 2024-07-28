package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"ecom-backend/internal/validator"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ProductCategoryHandler struct {
	BaseHandler
	productCategorySvc *service.ProductCategoryService
}

func NewProductCategoryHandler(logger *jsonlog.Logger, productCategorySvc *service.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{BaseHandler: BaseHandler{logger: logger}, productCategorySvc: productCategorySvc}
}

func (h *ProductCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	input := struct {
		Name     string  `json:"name"`
		ParentId *string `json:"parent_id"`
	}{}

	err := h.ReadJSON(w, r, &input)

	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.Name != "", "name", "must be provided")

	if input.ParentId != nil {
		if *input.ParentId == "" {
			input.ParentId = nil
		} else {
			v.Check(validator.IsValidUUID(*input.ParentId), "parent_id", "must be a valid UUID")
		}
	}

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.Errors)
		return
	}

	productCategoryRecord, err := h.productCategorySvc.CreateProductCategory(r.Context(), input.Name, input.ParentId)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	h.WriteJson(w, http.StatusCreated, ResponseBody{Payload: Envelope{"category": productCategoryRecord}}, nil)
}

func (h *ProductCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productCategorySvc.GetAll(r.Context())

	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	h.WriteJson(w, http.StatusOK, ResponseBody{Payload: Envelope{"categories": categories}}, nil)
}

func (h *ProductCategoryHandler) DeleteById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	categoryId := ps.ByName("categoryId")

	if categoryId == "" {
		h.NotFoundResponse(w, r)
		return
	}

	err := h.productCategorySvc.MarkAsDeleted(r.Context(), categoryId)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			h.NotFoundResponse(w, r)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	h.WriteJson(w, http.StatusOK, Envelope{"success": true}, nil)
}

func (h *ProductCategoryHandler) UpdateById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	categoryId := ps.ByName("categoryId")

	if categoryId == "" {
		h.NotFoundResponse(w, r)
		return
	}

	var input service.UpdateProductCategoryInput

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

	record, err := h.productCategorySvc.UpdateById(r.Context(), categoryId, input)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidProductCategory):
			h.FailedValidationResponse(w, r, map[string]string{"parent_id": "invalid parent category"})
		default:
			h.ServerErrorResponse(w, r, err)

		}
		return
	}

	h.WriteJson(w, http.StatusOK, ResponseBody{Payload: Envelope{"category": record}}, nil)
}

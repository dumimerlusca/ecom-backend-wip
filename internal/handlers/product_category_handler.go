package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/service"
	"ecom-backend/internal/validator"
	"net/http"
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

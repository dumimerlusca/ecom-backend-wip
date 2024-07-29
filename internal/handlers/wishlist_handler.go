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

type WishlistHandler struct {
	BaseHandler
	wishlistSvc *service.WishlistService
}

func NewWishlistHandler(logger *jsonlog.Logger, wishlistSvc *service.WishlistService) *WishlistHandler {
	return &WishlistHandler{BaseHandler: BaseHandler{logger: logger}, wishlistSvc: wishlistSvc}
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clientIdentifier := contextGetClientIdentifier(r)

	var input struct {
		VariantId string `json:"variant_id"`
	}

	err := h.ReadJSON(w, r, &input)

	if err != nil {
		h.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.VariantId != "", "variant_id", "must be provided")
	v.Check(validator.IsValidUUID(input.VariantId), "variant_id", "is not valid uuid")

	if !v.Valid() {
		h.FailedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = h.wishlistSvc.Insert(r.Context(), clientIdentifier, input.VariantId)

	if err != nil {
		if errors.Is(err, model.ErrProductAlreadyWishlisted) {
			h.BadRequestResponse(w, r, err)
			return
		}
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJson(w, http.StatusCreated, Envelope{"success": true}, nil)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}

}

func (h *WishlistHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdentifier := contextGetClientIdentifier(r)

	items, err := h.wishlistSvc.ListProducts(r.Context(), userIdentifier)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJson(w, http.StatusOK, ResponseBody{Payload: items}, nil)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

func (h *WishlistHandler) DeleteItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	itemId := ps.ByName("id")
	userIdentifier := contextGetClientIdentifier(r)

	if itemId == "" {
		h.NotFoundResponse(w, r)
		return
	}

	err := h.wishlistSvc.DeleteItem(r.Context(), userIdentifier, itemId)

	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			h.NotFoundResponse(w, r)
			return
		}

		if errors.Is(err, service.ErrUnauthorizedRequest) {
			h.UnauthorizedResponse(w, r)
			return
		}

		h.ServerErrorResponse(w, r, err)
		return
	}

	err = h.WriteJson(w, http.StatusOK, Envelope{"success": true}, nil)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

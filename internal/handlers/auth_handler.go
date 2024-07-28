package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"ecom-backend/internal/validator"
	"errors"
	"net/http"
)

type AuthHandler struct {
	BaseHandler
	authService *service.AuthService
}

func NewAuthHandler(logger *jsonlog.Logger, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{BaseHandler: BaseHandler{logger: logger}, authService: authService}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterUserInput

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

	user, err := h.authService.RegisterUser(r.Context(), &input)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrDuplicatedEmail):
			h.FailedValidationResponse(w, r, map[string]string{"email": "address already taken"})
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	h.WriteJson(w, http.StatusCreated, Envelope{"user": user}, nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginUserInput

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

	token, err := h.authService.LoginUser(r.Context(), input)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			h.InvalidCredentialsResponse(w, r)
		case errors.Is(err, service.ErrAccountActivationRequired):
			h.AccountActivationRequiredResponse(w, r)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = h.WriteJson(w, http.StatusOK, Envelope{"authentication_token": token}, nil)

	if err != nil {
		h.ServerErrorResponse(w, r, err)
	}
}

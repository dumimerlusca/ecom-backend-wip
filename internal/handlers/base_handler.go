package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/response"
	"encoding/json"
	"net/http"
)

type BaseHandler struct {
	logger *jsonlog.Logger
}

func (h *BaseHandler) WriteJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(js); err != nil {
		h.logger.PrintError(err, nil)
		return err
	}

	return nil
}

func (h *BaseHandler) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	w.WriteHeader(status)
	env := response.Envelope{"error": message}
	err := h.WriteJson(w, status, env, nil)

	if err != nil {
		h.LogError(r, err)
		w.WriteHeader(500)
	}
}

func (h *BaseHandler) LogError(r *http.Request, err error) {
	h.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (h *BaseHandler) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.LogError(r, err)
	msg := "the server encountered a problem and could not process your request"
	h.ErrorResponse(w, r, http.StatusInternalServerError, msg)
}

func (h *BaseHandler) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"
	h.ErrorResponse(w, r, http.StatusNotFound, msg)
}

func (h *BaseHandler) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *BaseHandler) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	h.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

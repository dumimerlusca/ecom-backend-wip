package handlers

import (
	"ecom-backend/internal/jsonlog"
	"errors"
	"fmt"
	"io"
	"strings"

	"encoding/json"
	"net/http"
)

type ResponseBody struct {
	Payload  interface{} `json:"payload,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}
type Envelope map[string]interface{}
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

// readJSON decodes request Body into corresponding Go type. It triages for any potential errors
// and returns corresponding appropriate errors.
func (h *BaseHandler) ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB to prevent
	// any potential nefarious DoS attacks.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. So, if the JSON from the client includes any field which
	// cannot be mapped to the target destination, the decoder will return an error
	// instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body to the destination.
	err := dec.Decode(dst)
	if err != nil {
		// If there is an error during decoding, start the error triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Use the error.As() function to check whether the error has the type *json.SyntaxError.
		// If it does, then return a plain-english error message which includes the location
		// of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at (charcter %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax error in the JSON. So, we check for this using errors.Is() and return
		// a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors.
		// These occur when the JSON value is the wrong type for the target destination.
		// If the error relates to a specific field, then we include that in our error message
		// to make it easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We check
		// for this with errors.Is() and return a plain-english error message instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message.
		// Note, that there's an open issue at https://github.com/golang/go/issues/29035
		// regarding turning this into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request body exceeds 1MB in size then decode will now fail with the
		// error "http: request body too large". There is an open issue about turning
		// this into a distinct error type at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an error
		// to our handler. At the end of this chapter we'll talk about panicking
		// versus returning, and discuss why it's an appropriate thing to do in this specific
		// situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else, return the error message as-is.
		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value then this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body, and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (h *BaseHandler) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	w.WriteHeader(status)
	env := Envelope{"error": message}
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

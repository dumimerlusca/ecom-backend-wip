package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Middleware struct {
	BaseHandler
	services *service.Services
}

func NewMiddleware(logger *jsonlog.Logger, services *service.Services) *Middleware {
	return &Middleware{BaseHandler: BaseHandler{logger: logger}, services: services}
}

func (mid *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic as
		// Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the response. This
				// acts a trigger to make Go's HTTP server automatically close the current
				// connection after a response has been sent.
				w.Header().Set("Connection:", "close")
				// The value returned by recover() has the type interface{}, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using our
				// custom Logger type at the ERROR level and send the client a
				// 500 Internal Server Error response.
				mid.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

var AnonymousUser = &model.UserRecord{}

func isAnonymousUser(user *model.UserRecord) bool {
	return user == AnonymousUser
}

func (mid *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")

		if authorizationHeader == "" {
			r = contextSetUser(r, AnonymousUser)

			sid := getSessionId(r)
			r = contextSetClientIdentifier(r, sid)

			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authorizationHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			mid.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := parts[1]

		user, err := mid.services.Token.GetUserByToken(r.Context(), token, service.ScopeAuthentication)

		if err != nil {
			switch {
			case errors.Is(err, model.ErrRecordNotFound):
				mid.InvalidAuthenticationTokenResponse(w, r)
			default:
				mid.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = contextSetUser(r, user)
		r = contextSetClientIdentifier(r, user.Id)

		next.ServeHTTP(w, r)
	})
}

func (mid *Middleware) RequireActivation(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := contextGetUser(r)

		if isAnonymousUser(user) {
			mid.UnauthorizedResponse(w, r)
			return
		}

		if !user.Activated {
			mid.ForbiddenResponse(w, r)
			return
		}

		next(w, r, ps)
	}
}

func (mid *Middleware) RequireSessionOrUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		clientIdentifier := contextGetClientIdentifier(r)

		if clientIdentifier == "" {
			mid.SessionOrUserRequiredResponse(w, r)
			return
		}

		next(w, r, ps)
	}
}

func (mid *Middleware) AdminOnly(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user := contextGetUser(r)

		if isAnonymousUser(user) {
			mid.UnauthorizedResponse(w, r)
			return
		}

		if !user.IsAdmin {
			mid.ForbiddenResponse(w, r)
			return
		}

		next(w, r, ps)
	}
}

func (mid *Middleware) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		// Add the "Vary: Access-Control-Request-Method" header.
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")

		// TODO Add the config back
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			// Check if the request has the HTTP method OPTIONS and contains the
			// "Access-Control-Request-Method" header. If it does, then we treat
			// it as a preflight request.
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				// Set the necessary preflight response headers, as discussed
				// previously.
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Session-ID")
				// Write the headers along with a 200 OK status and return from
				// the middleware with no further action.
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

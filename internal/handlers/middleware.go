package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	fmt.Println(user == AnonymousUser)
	return user == AnonymousUser
}

func (mid *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")

		if authorizationHeader == "" {
			r = contextSetUser(r, AnonymousUser)
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

		next.ServeHTTP(w, r)
	})
}

func (mid *Middleware) RequireActivation(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := contextGetUser(r)

		if isAnonymousUser(user) {
			mid.UnauthorizedResponse(w, r)
			return
		}

		if !user.Activated {
			mid.ForbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

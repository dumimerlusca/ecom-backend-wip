package handlers

import (
	"context"
	"ecom-backend/internal/model"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")
const clientIdentifierContextKey = contextKey("client-identifier") // user id for registerd users OR session id for guests

func contextSetClientIdentifier(r *http.Request, clientIdentifier string) *http.Request {
	ctx := context.WithValue(r.Context(), clientIdentifierContextKey, clientIdentifier)

	return r.WithContext(ctx)
}

func contextGetClientIdentifier(r *http.Request) string {
	clientIdentifier, ok := r.Context().Value(clientIdentifierContextKey).(string)

	if !ok {
		panic("missing client identifier value in request context")
	}

	return clientIdentifier
}

func contextSetUser(r *http.Request, user *model.UserRecord) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)

	return r.WithContext(ctx)
}

func contextGetUser(r *http.Request) *model.UserRecord {
	user, ok := r.Context().Value(userContextKey).(*model.UserRecord)

	if !ok {
		panic("missing user  value in request context")
	}

	return user
}

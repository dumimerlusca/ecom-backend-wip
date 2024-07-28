package handlers

import (
	"context"
	"ecom-backend/internal/model"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

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

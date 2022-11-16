// Filename: test2/cmd/api/context.go
package main

import (
	"context"
	"net/http"

	"michaelgomez.net/internal/data"
)

// Defining the custom key
type contextKey string

// Making the user key
const userContextKey = contextKey("user")

// Method that adds the user to the context
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// Method to retrieve the user struct
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}

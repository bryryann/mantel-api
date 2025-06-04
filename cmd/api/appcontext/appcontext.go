package appcontext

import (
	"context"
	"net/http"

	"github.com/bryryann/mantel/backend/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

type Context struct{}

func (c *Context) SetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (c *Context) GetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

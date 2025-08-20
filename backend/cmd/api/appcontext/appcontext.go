package appcontext

import (
	"context"
	"net/http"

	"github.com/bryanznk/mantel/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

type contextKey string

const (
	userContextKey   = contextKey("user")
	paramsContextKey = contextKey("params")
)

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

func (c *Context) SetParams(r *http.Request, ps httprouter.Params) *http.Request {
	ctx := context.WithValue(r.Context(), paramsContextKey, ps)
	return r.WithContext(ctx)
}

func (c *Context) GetParams(r *http.Request) httprouter.Params {
	ps, ok := r.Context().Value(paramsContextKey).(httprouter.Params)
	if !ok {
		return nil
	}
	return ps
}

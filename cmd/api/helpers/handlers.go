package helpers

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// AdaptHttpRouterHandle converts a http.HandlerFunc to an httprouter.Handle
func AdaptHttpRouterHandle(handler http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		for _, param := range ps {
			ctx = context.WithValue(ctx, param.Key, param.Value)
		}
		r = r.WithContext(ctx)

		handler(w, r)
	}
}

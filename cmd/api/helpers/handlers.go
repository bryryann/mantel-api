package helpers

import (
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/julienschmidt/httprouter"
)

// AdaptHttpRouterHandle converts a http.HandlerFunc to an httprouter.Handle
func AdaptHttpRouterHandle(ctx *appcontext.Context, handler http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		r = ctx.SetParams(r, ps)
		handler(w, r)
	}
}

func AdaptHttpHandlerFunc(ctx *appcontext.Context, handler httprouter.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ps := ctx.GetParams(r)
		handler(w, r, ps)
	}
}

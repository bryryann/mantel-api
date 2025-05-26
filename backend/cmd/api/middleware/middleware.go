package middleware

import (
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/bryryann/mantel/backend/internal/data"
)

func Apply(ctx *appcontext.Context, models *data.UserModel, router http.Handler) http.Handler {
	return authenticate(ctx, models, router)
}

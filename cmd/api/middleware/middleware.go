// Package middleware provides HTTP middleware utilities for the application.
package middleware

import (
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/bryryann/mantel/backend/internal/data"
)

// Apply wraps the provided HTTP handler with middleware for authentication.
// It takes the application context, app models, and the original HTTP handler as arguments.
// Returns a new HTTP handler that includes authentication middleware.
func Apply(ctx *appcontext.Context, models *data.Models, router http.Handler) http.Handler {
	return authenticate(ctx, &models.Users, router)
}

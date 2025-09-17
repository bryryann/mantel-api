// Package middleware provides HTTP middleware utilities for the application.
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/bryryann/mantel/backend/cmd/api/config"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/pascaldekloe/jwt"
)

// Apply wraps the provided HTTP handler with middleware for authentication.
// It takes the application context, app models, and the original HTTP handler as arguments.
// Returns a new HTTP handler that includes authentication middleware.
func Apply(ctx *appcontext.Context, models *data.Models, router http.Handler) http.Handler {
	return authenticate(ctx, &models.Users, router)
}

func ValidateToken(token string, cfg *config.Configuration) (*jwt.Claims, error) {
	claims, err := jwt.HMACCheck([]byte(token), []byte(cfg.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("invalid token signature: %w", err)
	}

	if !claims.Valid(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	if claims.Issuer != cfg.JWT.Issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	if cfg.Env != "DEVELOPMENT" {
		switch {
		case !claims.AcceptAudience(cfg.JWT.Audience):
			return nil, fmt.Errorf("invalid token audience")
		case claims.Issuer != cfg.JWT.Issuer:
			return nil, fmt.Errorf("invalid token issuer")
		}
	}

	return claims, nil
}

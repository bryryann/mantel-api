package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/bryryann/mantel/backend/cmd/api/config"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/pascaldekloe/jwt"
)

// authenticate is a middleware function that validates the JWT token from the Authorization header.
// It sets the user in the request context and passes the request to the next handler.
func authenticate(ctx *appcontext.Context, models *data.UserModel, next http.Handler) http.Handler {
	res := responses.Get()
	config := config.Load()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add Vary header for caching purposes
		w.Header().Add("Vary", "Authorization")

		// Retrieve the Authorization header
		authHeader := r.Header.Get("Authorization")

		// If no Authorization header is present, set the user as anonymous and proceed
		if authHeader == "" {
			r = ctx.SetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// Split the Authorization header into parts
		headerParts := strings.Split(authHeader, " ")
		// Validate the format of the Authorization header
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		// Extract the token from the header
		token := headerParts[1]

		// Verify the token using HMAC
		claims, err := jwt.HMACCheck([]byte(token), []byte(config.JWT.Secret))
		if err != nil {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check if the token is valid based on its expiration time
		if !claims.Valid(time.Now()) {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		// Verify the token issuer
		if claims.Issuer != config.JWT.Issuer {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		// Verify the token audience
		if claims.AcceptAudience(config.JWT.Audience) {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		// Parse the user ID from the token subject
		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			res.ServerErrorResponse(w, r, err)
			return
		}

		// Retrieve the user from the database using the user ID
		user, err := models.Get(userID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				res.InvalidAuthenticationTokenResponse(w, r)
			default:
				res.ServerErrorResponse(w, r, err)
			}
			return
		}

		// Set the authenticated user in the request context
		r = ctx.SetUser(r, user)

		// Pass the request to the next handler
		next.ServeHTTP(w, r)

	})
}

// RequireAuthenticatedUser is a middleware function that ensures the user is authenticated.
// If the user is anonymous, it responds with an authentication required error.
func RequireAuthenticatedUser(ctx *appcontext.Context, next http.HandlerFunc) http.HandlerFunc {
	res := responses.Get()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user from the request context
		user := ctx.GetUser(r)

		// If the user is anonymous, respond with an authentication required error
		if user.IsAnonymous() {
			res.AuthenticationRequiredResponse(w, r)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

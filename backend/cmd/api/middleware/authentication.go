package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/bryryann/mantel/backend/cmd/api/config"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
)

// authenticate is a middleware function that validates the JWT token from the Authorization header.
// It sets the user in the request context and passes the request to the next handler.
func authenticate(ctx *appcontext.Context, models *data.UserModel, next http.Handler) http.Handler {
	res := responses.Get()
	config := config.Load()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add Vary header for caching purposes
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = ctx.SetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		claims, err := ValidateToken(token, config)
		if err != nil {
			res.InvalidAuthenticationTokenResponse(w, r)
			fmt.Println("token validation error: ", err)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			res.ServerErrorResponse(w, r, err)
			return
		}

		user, err := models.Get(userID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				res.InvalidAuthenticationTokenResponse(w, r)
				fmt.Println("user not on db")
			default:
				res.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = ctx.SetUser(r, user)
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

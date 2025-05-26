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

func authenticate(ctx *appcontext.Context, models *data.UserModel, next http.Handler) http.Handler {
	res := responses.Get()
	config := config.Load()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		claims, err := jwt.HMACCheck([]byte(token), []byte(config.JWT.Secret))
		if err != nil {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.Valid(time.Now()) {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		if claims.Issuer != config.JWT.Issuer {
			res.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		if claims.AcceptAudience(config.JWT.Audience) {
			res.InvalidAuthenticationTokenResponse(w, r)
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
			default:
				res.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = ctx.SetUser(r, user)

		next.ServeHTTP(w, r)

	})
}

func RequireAuthenticatedUser(ctx *appcontext.Context, next http.HandlerFunc) http.HandlerFunc {
	res := responses.Get()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ctx.GetUser(r)

		if user.IsAnonymous() {
			res.AuthenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

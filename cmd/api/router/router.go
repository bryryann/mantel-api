// Package handlers groups every API handler and initializes them.
package router

import (
	"net/http"
	"sync"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/middleware"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

// Route represents a registered HTTP route, containing it's handler and method.
type Route struct {
	Path    string            // URL Path pattern
	Method  string            // HTTP Method(GET, POST, etc)
	Handler httprouter.Handle // Handler function for this route.
}

type envelope map[string]any

var (
	routes []Route
	mu     sync.RWMutex
)

// SetupRouter initializes an http.Handler with all registered routes.
func SetupRouter(ctx *appcontext.Context, models *data.Models) http.Handler {
	router := httprouter.New()
	mu.RLock()
	defer mu.RUnlock()

	// TODO: Add NotFound and MethodNotAllowed handlers.

	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handler)
	}

	return middleware.Apply(ctx, models, router)
}

func init() {
	Get("/v1/healthcheck", helpers.AdaptHttpRouterHandle(healthCheck))

	Post("/v1/users", helpers.AdaptHttpRouterHandle(registerUser))
	Post("/v1/users/:follower_id/follow", followUser)
	Delete("/v1/users/:follower_id/unfollow/:followee_id", unfollowUser)

	Get("/v1/users/:user_id/followers", userFollowers)
	Get("/v1/users/:user_id/followees", userFollowees)

	Post("/v1/tokens/authentication", helpers.AdaptHttpRouterHandle(authenticateToken))
}

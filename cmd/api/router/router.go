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

func InitializeRouter(ctx *appcontext.Context) {
	httprouterCompatible := helpers.AdaptHttpRouterHandle
	httpCompatible := helpers.AdaptHttpHandlerFunc

	// healthcheck
	Get("/v1/healthcheck", httprouterCompatible(ctx, healthCheck))

	// user and authentication
	ProtectedGet("/v1/users", getAuthUser, ctx)
	Get("/v1/users/:user_id", getUserByID)
	Post("/v1/users", httprouterCompatible(ctx, registerUser))
	Post("/v1/tokens/authentication", httprouterCompatible(ctx, authenticateToken))

	// follows
	Get("/v1/users/:user_id/followers", listUserFollowers)
	Get("/v1/users/:user_id/followees", listUserFollowees)
	ProtectedPost("/v1/users/:follower_id/follow", httpCompatible(ctx, followUser), ctx)
	ProtectedPost("/v1/users/:follower_id/unfollow/:followee_id", httpCompatible(ctx, unfollowUser), ctx)

	// friendships
	ProtectedPost("/v1/friends/send", sendFriendRequest, ctx)
	ProtectedPost("/v1/friends/accept", acceptFriendRequest, ctx)
}

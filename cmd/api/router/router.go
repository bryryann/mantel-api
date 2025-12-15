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
	Get("/v1/users/:user_id/follows/:followee_id", checkFollowStatus)
	ProtectedPost("/v1/users/:follower_id/follow", httpCompatible(ctx, followUser), ctx)
	ProtectedPost("/v1/users/:follower_id/unfollow/:followee_id", httpCompatible(ctx, unfollowUser), ctx)

	// friendships
	ProtectedGet("/v1/friend-requests", listPendingRequests, ctx)
	ProtectedPost("/v1/friend-requests", sendFriendRequest, ctx)
	ProtectedDelete("/v1/friend-requests/:request_id/reject", httpCompatible(ctx, rejectFriendRequest), ctx)
	ProtectedPatch("/v1/friend-requests/:id", httpCompatible(ctx, patchPendingFriendRequest), ctx)

	ProtectedGet("/v1/users/:user_id/friends", httpCompatible(ctx, getFriendsById), ctx)
	Get("/v1/users/:user_id/friends/:friend_id", getFriendshipStatus)

	// posts
	Get("/v1/posts/:post_id", findPostByID)
	ProtectedPost("/v1/posts", createNewPost, ctx)
	ProtectedDelete("/v1/posts/:post_id", httpCompatible(ctx, deletePostFromAuthUser), ctx)
	ProtectedPatch("/v1/posts/:post_id", httpCompatible(ctx, editPostContent), ctx)
	Get("/v1/users/:user_id/posts", getPostsFromUser)
	Get("/v1/users/:user_id/posts/:post_id", findPostByIDFromUser)

	// likes
	ProtectedPost("/v1/posts/:post_id/likes", httpCompatible(ctx, likePost), ctx)
	ProtectedDelete("/v1/posts/:post_id/likes", httpCompatible(ctx, dislikePost), ctx)
	Get("/v1/posts/:post_id/likes", listLikesOnPost)
	Get("/v1/posts/:post_id/likes/count", countLikesFromPost)
	Get("/v1/users/:user_id/liked/:post_id", hasUserLiked)
}

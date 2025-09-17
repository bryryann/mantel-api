package router

import (
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/middleware"
	"github.com/julienschmidt/httprouter"
)

// RegisterHandler adds a new HTTP route with explicit method specification
//
// Parameters:
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
//   - path: URL path pattern
//   - handler: Function to handle requests
func RegisterHandler(method, path string, handler httprouter.Handle) {
	mu.Lock()
	defer mu.Unlock()
	routes = append(routes, Route{
		Path:    path,
		Method:  method,
		Handler: handler,
	})
}

// Get register a handler for HTTP GET requests.
// This is a convenience wrapper aruond RegisterHandler.
func Get(path string, handler httprouter.Handle) {
	RegisterHandler(http.MethodGet, path, handler)
}

// ProtectedGet register a handler for HTTP GET requests that require an authenticated user.
func ProtectedGet(path string, handler http.HandlerFunc, ctx *appcontext.Context) {
	protectedHandler := helpers.AdaptHttpRouterHandle(ctx, middleware.RequireAuthenticatedUser(ctx, handler))
	Get(path, protectedHandler)
}

// Post register a handler for HTTP POST requests.
// This is a convenience wrapper aruond RegisterHandler.
func Post(path string, handler httprouter.Handle) {
	RegisterHandler(http.MethodPost, path, handler)
}

// ProtectedPost register a handler for HTTP POST requests that require an authenticated user.
func ProtectedPost(path string, handler http.HandlerFunc, ctx *appcontext.Context) {
	protectedHandler := helpers.AdaptHttpRouterHandle(ctx, middleware.RequireAuthenticatedUser(ctx, handler))
	Post(path, protectedHandler)
}

// Put register a handler for HTTP PUT requests.
// This is a convenience wrapper aruond RegisterHandler.
func Put(path string, handler httprouter.Handle) {
	RegisterHandler(http.MethodPut, path, handler)
}

// ProtectedPut register a handler for HTTP PUT requests that require an authenticated user.
func ProtectedPut(path string, handler http.HandlerFunc, ctx *appcontext.Context) {
	protectedHandler := helpers.AdaptHttpRouterHandle(ctx, middleware.RequireAuthenticatedUser(ctx, handler))
	Put(path, protectedHandler)
}

func Patch(path string, handler httprouter.Handle) {
	RegisterHandler(http.MethodPatch, path, handler)
}

func ProtectedPatch(path string, handler http.HandlerFunc, ctx *appcontext.Context) {
	protectedHandler := helpers.AdaptHttpRouterHandle(ctx, middleware.RequireAuthenticatedUser(ctx, handler))
	Patch(path, protectedHandler)
}

// Delete register a handler for HTTP DELETE requests.
// This is a convenience wrapper aruond RegisterHandler.
func Delete(path string, handler httprouter.Handle) {
	RegisterHandler(http.MethodDelete, path, handler)
}

// ProtectedDelete register a handler for HTTP DELETE requests that require an authenticated user.
func ProtectedDelete(path string, handler http.HandlerFunc, ctx *appcontext.Context) {
	protectedHandler := helpers.AdaptHttpRouterHandle(ctx, middleware.RequireAuthenticatedUser(ctx, handler))
	Delete(path, protectedHandler)
}

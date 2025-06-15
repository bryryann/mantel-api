package router

import (
	"net/http"

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

// Post register a handler for HTTP POST requests.
// This is a convenience wrapper aruond RegisterHandler.
func Post(path string, handler httprouter.Handle) {
	RegisterHandler(http.MethodPost, path, handler)
}

// Put register a handler for HTTP PUT requests.
// This is a convenience wrapper aruond RegisterHandler.
func Put(path string, handler httprouter.Handle) {
	RegisterHandler(http.MethodPut, path, handler)
}

// Delete register a handler for HTTP DELETE requests.
// This is a convenience wrapper aruond RegisterHandler.
func Delete(path string, handler httprouter.Handle) {
	RegisterHandler(http.MethodDelete, path, handler)
}

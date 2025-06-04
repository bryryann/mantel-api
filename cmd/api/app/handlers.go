package app

import "net/http"

// RegisterHandler adds a new HTTP route with explicit method specification
//
// Parameters:
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
//   - path: URL path pattern
//   - handler: Function to handle requests
func (a *App) RegisterHandler(method, path string, handler http.HandlerFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.routes = append(a.routes, Route{
		Path:    path,
		Method:  method,
		Handler: handler,
	})
}

// Get register a handler for HTTP GET requests.
// This is a convenience wrapper aruond RegisterHandler.
func (a *App) Get(path string, handler http.HandlerFunc) {
	a.RegisterHandler(http.MethodGet, path, handler)
}

// Post register a handler for HTTP POST requests.
// This is a convenience wrapper aruond RegisterHandler.
func (a *App) Post(path string, handler http.HandlerFunc) {
	a.RegisterHandler(http.MethodPost, path, handler)
}

// Put register a handler for HTTP PUT requests.
// This is a convenience wrapper aruond RegisterHandler.
func (a *App) Put(path string, handler http.HandlerFunc) {
	a.RegisterHandler(http.MethodPut, path, handler)
}

// Delete register a handler for HTTP DELETE requests.
// This is a convenience wrapper aruond RegisterHandler.
func (a *App) Delete(path string, handler http.HandlerFunc) {
	a.RegisterHandler(http.MethodDelete, path, handler)
}

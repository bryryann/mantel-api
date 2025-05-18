// Package app provides the application dependency injection container
// and configuration registration system.
package app

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Route represents a registered HTTP route, containing it's handler and method.
type Route struct {
	Path    string           // URL Path pattern
	Method  string           // HTTP Method(GET, POST, etc)
	Handler http.HandlerFunc // Handler function for this route.
}

// App is the application container that holds:
// - Shared dependencies
// - Registered HTTP routes
// - Thread-safe synchronization
type App struct {
	routes []Route
	mu     sync.RWMutex
}

var (
	instance *App      // Singleton instance.
	once     sync.Once // Ensures singleton initialization happens once.
)

// Get returns the singleton App instance, initializing it if necessary.
// Concurrent-safe.
func Get() *App {
	once.Do(func() {
		instance = &App{
			routes: make([]Route, 0),
		}
	})

	return instance
}

// SetupRouter initializes an http.Handler with all registered routes.
func (a *App) SetupRouter() http.Handler {
	router := httprouter.New()
	a.mu.RLock()
	defer a.mu.RUnlock()

	// TODO: Add NotFound and MethodNotAllowed handlers.

	for _, route := range a.routes {
		router.HandlerFunc(route.Method, route.Path, route.Handler)
	}
	return router
}

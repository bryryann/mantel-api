package app

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

type App struct {
	routes []Route
	mu     sync.RWMutex
}

var (
	instance *App
	once     sync.Once
)

func Get() *App {
	once.Do(func() {
		instance = &App{
			routes: make([]Route, 0),
		}
	})

	return instance
}

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

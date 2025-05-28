// Package app provides the application dependency injection container
// and configuration registration system.
package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/bryryann/mantel/backend/cmd/api/appcontext"
	"github.com/bryryann/mantel/backend/cmd/api/config"
	"github.com/bryryann/mantel/backend/cmd/api/database"
	"github.com/bryryann/mantel/backend/cmd/api/middleware"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

// Route represents a registered HTTP route, containing it's handler and method.
type Route struct {
	Path    string           // URL Path pattern
	Method  string           // HTTP Method(GET, POST, etc)
	Handler http.HandlerFunc // Handler function for this route.
}

// App is the application container that holds:
// - Shared dependencies (Configuration, Logger)
// - Registered HTTP routes
// - Thread-safe synchronization
type App struct {
	Config    *config.Configuration
	Database  *database.Database
	Logger    *slog.Logger
	Models    *data.Models
	Context   *appcontext.Context
	Responses *responses.Responses
	routes    []Route
	mu        sync.RWMutex
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
			routes:    make([]Route, 0),
			Context:   &appcontext.Context{},
			Responses: responses.Get(),
		}
	})

	return instance
}

// SetConfig attributes a *config.Configuration as the app Config.
func (a *App) SetConfig(cfg *config.Configuration) {
	a.Config = cfg
}

// SetDB opens a sql connection and attributes a *database.Database to the app.
func (a *App) SetDB(dsn string) error {
	db, err := database.OpenConnection(dsn)
	if err != nil {
		return fmt.Errorf("failed to create db connection: %w", err)
	}

	a.Database = db
	return nil
}

func (a *App) SetModels() {
	if a.Database == nil {
		a.Logger.Error("failed to set models. no proper database")
		return
	}

	a.Models = data.NewModels(a.Database.DB)
}

// ConfigureLogger sets the global application logger.
func (a *App) ConfigureLogger(logLevel string) {
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)

	a.Logger = slog.New(handler)
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
	return middleware.Apply(a.Context, a.Models, router)
}

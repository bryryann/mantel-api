package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/config"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/router"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const ascii = `
                       __         .__
  _____ _____    _____/  |_  ____ |  |
 /     \\__  \  /    \   __\/ __ \|  |
|  Y Y  \/ __ \|   |  \  | \  ___/|  |__
|__|_|__(_____/|___|  /__|  \_____>____/ `

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func main() {
	fmt.Println(ascii)

	dsn := helpers.GetEnvString("MANTEL_DB_DSN", "")
	if dsn == "" {
		panic("invalid db dsn\n")
	}
	application := app.Get()

	cfg := config.Load()

	application.ConfigureLogger("info")

	application.SetConfig(cfg)
	application.SetDB(cfg.DSN)
	application.SetModels()

	router.InitializeRouter(application.Context)

	application.Logger.Info("all set up!")

	startServer()
}

// startServer contains all code related to api initialization.
func startServer() {
	app := app.Get()
	router := router.SetupRouter(app.Context, app.Models)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
	}

	err := srv.ListenAndServe()
	if err != nil {
		app.Logger.Error(err.Error())
	}
}

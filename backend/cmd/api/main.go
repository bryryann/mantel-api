package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/config"
	_ "github.com/bryryann/mantel/backend/cmd/api/handlers"
	"github.com/joho/godotenv"
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

	cfg := config.Load()

	app.Get().SetConfig(cfg)

	startServer()
}

func startServer() {
	application := app.Get()
	router := application.SetupRouter()

	// TODO: Add logger
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", application.Config.Port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("fail: %v", err)
	}
}

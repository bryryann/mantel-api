package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	_ "github.com/bryryann/mantel/backend/cmd/api/handlers"
)

const ascii = `
                       __         .__
  _____ _____    _____/  |_  ____ |  |
 /     \\__  \  /    \   __\/ __ \|  |
|  Y Y  \/ __ \|   |  \  | \  ___/|  |__
|__|_|__(_____/|___|  /__|  \_____>____/ `

func main() {
	fmt.Println(ascii)

	startServer()
}

func startServer() {
	application := app.Get()
	router := application.SetupRouter()

	// TODO: Add logger
	// TODO: Import address and other sensitive info to env variables.
	srv := &http.Server{
		Addr:         ":4000",
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(fmt.Errorf("fail: %v", err))
	}
}

package handlers

import (
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/app"
)

func init() {
	app := app.Get()

	app.RegisterHandler(http.MethodGet, "/v1/healthcheck", HealthCheck)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Mantel version 0.0.1"))
}

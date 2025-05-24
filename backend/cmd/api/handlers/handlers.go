// Package handlers groups every API handler and initializes them.
package handlers

import "github.com/bryryann/mantel/backend/cmd/api/app"

func init() {
	app := app.Get()

	app.Get("/v1/healthcheck", HealthCheck)

	app.Post("/v1/users", registerUserHandler)
}

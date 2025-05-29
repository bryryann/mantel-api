// Package handlers groups every API handler and initializes them.
package handlers

import "github.com/bryryann/mantel/backend/cmd/api/app"

func init() {
	app := app.Get()

	app.Get("/v1/healthcheck", healthCheck)

	app.Post("/v1/users", registerUser)
	app.Post("/v1/tokens/authentication", authTokenHandler)
}

// Package handlers groups every API handler and initializes them.
package handlers

import (
	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
)

type envelope map[string]any

func init() {
	app := app.Get()

	app.Get("/v1/healthcheck", helpers.AdaptHttpRouterHandle(healthCheck))

	app.Post("/v1/users", helpers.AdaptHttpRouterHandle(registerUser))
	app.Post("/v1/users/:follower_id/follow", followUser)
	app.Delete("/v1/users/:follower_id/unfollow/:followee_id", unfollowUser)

	app.Post("/v1/tokens/authentication", helpers.AdaptHttpRouterHandle(authenticateToken))
}

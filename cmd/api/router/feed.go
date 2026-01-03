package router

import (
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/jsonhttp"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
)

func getFeed(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	query := r.URL.Query()
	page := helpers.ParseIntOrDefault(query.Get("page"), 1)
	pageSize := helpers.ParseIntOrDefault(query.Get("page_size"), 20)

	paginationData := data.Pagination{
		Page:     page,
		PageSize: pageSize,
		Sort:     "",
	}
	posts, err := app.Models.Feed.Fetch(user.ID, paginationData)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if posts == nil {
		posts = []data.PostPublic{}
	}

	jsonResponse := envelope{
		"feed": posts,
		"meta": map[string]any{
			"page":      page,
			"page_size": pageSize,
		},
	}
	err = jsonhttp.WriteJSON(w, http.StatusOK, jsonResponse, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

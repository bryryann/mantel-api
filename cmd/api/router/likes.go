package router

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/jsonhttp"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

func likePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	postIDParam := ps.ByName("post_id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	user := app.Context.GetUser(r)

	exists, err := app.Models.Posts.Exists(int64(postID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if !exists {
		res.NotFoundResponse(w, r)
		return
	}

	like, err := app.Models.Likes.Like(int64(user.ID), int64(postID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if like == nil {
		jsonResponse := envelope{
			"message": "post already liked",
		}

		err = jsonhttp.WriteJSON(w, http.StatusOK, jsonResponse, nil)
		if err != nil {
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	jsonResponse := envelope{
		"like": like,
	}
	err = jsonhttp.WriteJSON(w, http.StatusOK, jsonResponse, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func dislikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	postIDParam := ps.ByName("post_id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	user := app.Context.GetUser(r)

	exists, err := app.Models.Posts.Exists(int64(postID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if !exists {
		res.NotFoundResponse(w, r)
		return
	}

	err = app.Models.Likes.Dislike(user.ID, int64(postID))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			res.NotFoundResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhttp.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func listLikesOnPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()
	query := r.URL.Query()

	postIDParam := ps.ByName("post_id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	exists, err := app.Models.Posts.Exists(int64(postID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if !exists {
		res.NotFoundResponse(w, r)
		return
	}

	page := helpers.ParseIntOrDefault(query.Get("page"), 1)
	pageSize := helpers.ParseIntOrDefault(query.Get("page_size"), 20)
	sort := query.Get("sort")

	paginationData := data.Pagination{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	}

	likes, err := app.Models.Likes.ListLikesFromPost(int64(postID), paginationData)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if likes == nil {
		likes = []data.LikePublic{}
	}

	jsonResponse := envelope{
		"likes": likes,
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

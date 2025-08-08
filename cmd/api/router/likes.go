package router

import (
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/jsonhttp"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
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

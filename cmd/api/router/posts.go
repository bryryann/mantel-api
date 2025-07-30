package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/jsonhttp"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

func findPostByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	postIDQuery := ps.ByName("post_id")

	postID, err := strconv.Atoi(postIDQuery)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	post, err := app.Models.Posts.Get(int64(postID))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			res.NotFoundResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	post.ID = int64(postID)

	err = jsonhttp.WriteJSON(w, http.StatusAccepted, envelope{"post": post}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func createNewPost(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	var input struct {
		Content string `json:"content"`
	}

	err := jsonhttp.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	post := &data.Post{
		UserID:  user.ID,
		Content: input.Content,
	}

	err = app.Models.Posts.Insert(post)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	err = jsonhttp.WriteJSON(w, http.StatusAccepted, envelope{"post": post}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

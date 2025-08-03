package router

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/jsonhttp"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/bryryann/mantel/backend/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func findPostByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	postIDParam := ps.ByName("post_id")

	postID, err := strconv.Atoi(postIDParam)
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

func findPostByIDFromUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	userID, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	postID, err := strconv.Atoi(ps.ByName("post_id"))
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	post, err := app.Models.Posts.FindByIDFromUser(int64(postID), int64(userID))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			message := fmt.Sprintf("post %d does not belong to user %d", postID, userID)
			res.ErrorResponse(w, r, http.StatusNotFound, message)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhttp.WriteJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func getPostsFromUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	userIDParam := ps.ByName("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	posts, err := app.Models.Posts.SelectAllFromUser(int64(userID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if posts == nil {
		posts = []data.PostPublic{}
	}

	err = jsonhttp.WriteJSON(w, http.StatusAccepted, envelope{"posts": posts}, nil)
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

	v := validator.New()
	if data.ValidatePost(v, post); !v.Valid() {
		res.FailedValidationResponse(w, r, v.Errors)
		return
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

func deletePostFromAuthUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	postIDParam := ps.ByName("post_id")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	exists, err := app.Models.Posts.CheckPostOwnership(int64(postID), user.ID)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if !exists {
		res.NotAuthorizedResponse(w, r)
		return
	}

	err = app.Models.Posts.Delete(int64(postID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	message := fmt.Sprintf("succesfully deleted post with id - %d", postID)
	err = jsonhttp.WriteJSON(w, http.StatusAccepted, message, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

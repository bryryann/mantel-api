package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

// followUser handles the request for following a user.
// It receives a follower_id in the url, and the person to be followed
// as followee_id as the JSON body.
func followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	idParam := ps.ByName("follower_id")
	followerID, err := strconv.Atoi(idParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	var input struct {
		FolloweeID int `json:"followee_id"`
	}

	err = helpers.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	// TODO: Add validation
	// - Verify whether follower_id/followee_id actually represents a registered user
	// ...
	// ...
	// ...

	err = app.Models.Follows.Insert(int64(followerID), int64(input.FolloweeID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(
		w,
		http.StatusCreated,
		envelope{"follower_id": followerID, "followee_id": input.FolloweeID},
		nil,
	)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

// unfollowUser deletes a follow instance from the database.
// It receives both the follower_id and followee_id, validates whether it exists or not, and perform the appropriate db query.
func unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	followerIDParam := ps.ByName("follower_id")
	followeeIDParam := ps.ByName("followee_id")

	followerID, err := strconv.Atoi(followerIDParam)
	followeeID, err := strconv.Atoi(followeeIDParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	err = app.Models.Follows.Delete(int64(followerID), int64(followeeID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func listUserFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	id, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	_, err = app.Models.Users.Exists(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUserNotFound):
			res.NotFoundResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	followers, err := app.Models.Follows.GetFollowers(int64(id))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if followers == nil {
		followers = []data.UserPublic{}
	}

	helpers.WriteJSON(w, http.StatusAccepted, envelope{"followers": followers}, nil)
}

func listUserFollowees(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	id, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	_, err = app.Models.Users.Exists(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUserNotFound):
			res.NotFoundResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	followees, err := app.Models.Follows.GetFollowees(int64(id))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if followees == nil {
		followees = []data.UserPublic{}
	}

	helpers.WriteJSON(w, http.StatusAccepted, envelope{"followees": followees}, nil)
}

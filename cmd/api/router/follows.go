package router

import (
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/julienschmidt/httprouter"
)

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

func userFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	id, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	user, err := app.Models.Users.Get(int64(id))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	followers, err := app.Models.Follows.GetFollowers(int64(id))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	env := envelope{"user": user, "followers": followers}

	helpers.WriteJSON(w, http.StatusAccepted, env, nil)
}

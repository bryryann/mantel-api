package router

import (
	"errors"
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
)

func sendFriendRequest(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	var input struct {
		UserID   int `json:"user_id"`
		FriendID int `json:"friend_id"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	fs := &data.Friendship{
		UserID:   int64(input.UserID),
		FriendID: int64(input.FriendID),
	}

	err = app.Models.Friendships.SendRequest(fs)
	if err != nil {
		switch err {
		case data.ErrFriendshipRequestToSelf:
			res.BadRequestResponse(w, r, err)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message":    "friend request sent",
		"created_at": fs.CreatedAt,
		"status":     fs.Status,
	}
	err = helpers.WriteJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func acceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	var input struct {
		UserID   int `json:"user_id"`
		FriendID int `json:"friend_id"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	fs := &data.Friendship{
		UserID:   int64(input.UserID),
		FriendID: int64(input.FriendID),
	}

	err = app.Models.Friendships.AcceptRequest(fs)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoSuchRequest):
			res.NotFoundResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = helpers.WriteJSON(w, http.StatusAccepted, envelope{"message": "friend request accepted"}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

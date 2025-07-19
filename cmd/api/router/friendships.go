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

	user := app.Context.GetUser(r)

	var input struct {
		ReceiverID int `json:"receiver_id"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	fs := &data.Friendship{
		SenderID:   user.ID,
		ReceiverID: int64(input.ReceiverID),
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

func acceptPendingFriendRequest(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	var input struct {
		SenderID int `json:"sender_id"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	fs := &data.Friendship{
		SenderID:   int64(input.SenderID),
		ReceiverID: user.ID,
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

func rejectPendingFriendRequest(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	var input struct {
		SenderID int `json:"sender_id"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	fs := &data.Friendship{
		SenderID:   int64(input.SenderID),
		ReceiverID: user.ID,
	}

	err = app.Models.Friendships.RejectRequest(fs)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoSuchRequest):
			res.NotFoundResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = helpers.WriteJSON(w, http.StatusAccepted, envelope{"message": "friend request rejected"}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func listPendingRequests(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	reqs, err := app.Models.Friendships.GetPendingRequests(int64(user.ID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusAccepted, envelope{"requests": reqs}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}
